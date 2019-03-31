package nosql

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
)

const sectionName = "mongodb"
const schemePrefix = "mongodb://"

// IniFile is interface to read INI file variables
type IniFile interface {
	Get(section string, name string) string
}

// Databases represents mongodb database pool
type Databases struct {
	ini  IniFile
	dbs  map[string]*Database
	opts []*options.ClientOptions
}

// Open returns mongodb database pool
func Open(ini IniFile, opts ...*options.ClientOptions) *Databases {
	return &Databases{
		ini:  ini,
		dbs:  map[string]*Database{},
		opts: opts,
	}
}

// Close closes mongo database pool
func (d *Databases) Close(ctx context.Context) error {
	for _, db := range d.dbs {
		if err := db.Client().Disconnect(ctx); err != nil {
			return err
		}
	}
	return nil
}

// Get returns database hanle
func (d *Databases) Get(ctx context.Context, name string,
	dbopts ...*options.DatabaseOptions) (*Database, error) {
	if db, ok := d.dbs[name]; ok {
		return db, nil
	}
	uri, dbname, err := d.getURI(name)
	if err != nil {
		return nil, err
	}
	opt := options.Client().
		SetReadConcern(readconcern.Majority()).
		ApplyURI(uri)
	opts := append([]*options.ClientOptions{opt}, d.opts...)
	client, err := mongo.Connect(ctx, options.MergeClientOptions(opts...))
	if err != nil {
		return nil, err
	}
	db := NewDatabase(client.Database(dbname, dbopts...))
	d.dbs[name] = db
	return db, nil
}

// getURI returns database URI and name using INI file substitution
func (d *Databases) getURI(name string) (string, string, error) {
	uri := d.ini.Get(sectionName, name)
	if uri == "" {
		return "", "", fmt.Errorf("variable %s not found in section %s", name, sectionName)
	}
	if !strings.HasPrefix(uri, schemePrefix) {
		uri = schemePrefix + uri
	}
	u, err := url.Parse(uri)
	if err != nil {
		return "", "", err
	}
	// remove dbname
	dbname := strings.TrimPrefix(u.Path, "/")
	u.Path = "/"
	// user = username:password or separate lines, user = username and password = password
	users := strings.SplitN(d.ini.Get(u.Host, "user"), ":", 2)
	if len(users) == 2 {
		u.User = url.UserPassword(users[0], users[1])
	} else if len(users) == 1 {
		if pass := d.ini.Get(u.Host, "password"); pass != "" {
			u.User = url.UserPassword(users[0], pass)
		}
	}
	// query = tag=value
	if query := d.ini.Get(u.Host, "query"); query != "" {
		if u.RawQuery != "" {
			u.RawQuery += "&"
		}
		u.RawQuery += query
	}
	// host substitution must be latest
	if host := d.ini.Get(u.Host, "host"); host != "" {
		u.Host = host
	}
	// converting each "%AB" into hex-decoded byte
	if uri, err = url.QueryUnescape(u.String()); err != nil {
		return "", "", err
	}
	return uri, dbname, nil
}
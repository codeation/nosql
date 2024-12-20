package nosql

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const sectionName = "mongodb"
const schemePrefix = "mongodb://"

var errorVariable = errors.New("variable not found")

const sides = 2

// IniFile is interface to read INI file variables.
type IniFile interface {
	Get(section string, name string) string
}

// Databases represents mongodb database pool.
type Databases struct {
	ini        IniFile
	dbs        map[string]*Database
	clientOpts []*options.ClientOptions
}

// Open returns mongodb database pool.
func Open(ini IniFile, opts ...*options.ClientOptions) *Databases {
	return &Databases{
		ini:        ini,
		dbs:        map[string]*Database{},
		clientOpts: opts,
	}
}

// Close closes mongo database pool.
func (d *Databases) Close(ctx context.Context) error {
	for _, db := range d.dbs {
		if err := db.Client().Disconnect(ctx); err != nil {
			return err
		}
	}

	return nil
}

// Get returns database hanle.
func (d *Databases) Get(ctx context.Context, name string, opts ...options.Lister[options.DatabaseOptions],
) (*Database, error) {
	if db, ok := d.dbs[name]; ok {
		return db, nil
	}

	uri, dbname, err := d.getURI(name)
	if err != nil {
		return nil, err
	}

	client, err := mongo.Connect(
		append([]*options.ClientOptions{options.Client().ApplyURI(uri)}, d.clientOpts...)...)
	if err != nil {
		return nil, err
	}

	db := NewDatabase(client.Database(dbname, opts...))
	d.dbs[name] = db

	return db, nil
}

// getURI returns database URI and name using INI file substitution.
func (d *Databases) getURI(name string) (string, string, error) {
	uri := d.ini.Get(sectionName, name)
	if uri == "" {
		return "", "", fmt.Errorf("%w ([%s] %s)", errorVariable, sectionName, name)
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
	if d.ini.Get(u.Host, "user") != "" && d.ini.Get(u.Host, "password") != "" {
		u.User = url.UserPassword(d.ini.Get(u.Host, "user"), d.ini.Get(u.Host, "password"))
	} else if users := strings.SplitN(d.ini.Get(u.Host, "user"), ":", sides); len(users) == sides {
		u.User = url.UserPassword(users[0], users[1])
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

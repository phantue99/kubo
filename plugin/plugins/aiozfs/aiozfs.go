package aiozfs

import (
	"fmt"
	"path/filepath"

	"github.com/ipfs/kubo/plugin"
	"github.com/ipfs/kubo/repo"
	"github.com/ipfs/kubo/repo/fsrepo"

	aiozfs "github.com/phantue99/go-ds-aiozfs"
)

// Plugins is exported list of plugins that will be loaded
var Plugins = []plugin.Plugin{
	&aiozfsPlugin{},
}

type aiozfsPlugin struct{}

var _ plugin.PluginDatastore = (*aiozfsPlugin)(nil)

func (*aiozfsPlugin) Name() string {
	return "ds-aiozfs"
}

func (*aiozfsPlugin) Version() string {
	return "0.1.0"
}

func (*aiozfsPlugin) Init(_ *plugin.Environment) error {
	return nil
}

func (*aiozfsPlugin) DatastoreTypeName() string {
	return "aiozfs"
}

type datastoreConfig struct {
	path      string
	shardFun  *aiozfs.ShardIdV1
	syncField bool
}

// BadgerdsDatastoreConfig returns a configuration stub for a badger datastore
// from the given parameters
func (*aiozfsPlugin) DatastoreConfigParser() fsrepo.ConfigFromMap {
	return func(params map[string]interface{}) (fsrepo.DatastoreConfig, error) {
		var c datastoreConfig
		var ok bool
		var err error

		c.path, ok = params["path"].(string)
		if !ok {
			return nil, fmt.Errorf("'path' field is missing or not boolean")
		}

		sshardFun, ok := params["shardFunc"].(string)
		if !ok {
			return nil, fmt.Errorf("'shardFunc' field is missing or not a string")
		}
		c.shardFun, err = aiozfs.ParseShardFunc(sshardFun)
		if err != nil {
			return nil, err
		}

		c.syncField, ok = params["sync"].(bool)
		if !ok {
			return nil, fmt.Errorf("'sync' field is missing or not boolean")
		}
		return &c, nil
	}
}

func (c *datastoreConfig) DiskSpec() fsrepo.DiskSpec {
	return map[string]interface{}{
		"type":      "aiozfs",
		"path":      c.path,
		"shardFunc": c.shardFun.String(),
	}
}

func (c *datastoreConfig) Create(path string) (repo.Datastore, error) {
	p := c.path
	if !filepath.IsAbs(p) {
		p = filepath.Join(path, p)
	}

	return aiozfs.CreateOrOpen(p, c.shardFun, c.syncField)
}

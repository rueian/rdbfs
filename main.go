package main

import (
	"os"

	"errors"

	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"
	"github.com/rueian/rdbfs/filesystem"
	"github.com/rueian/rdbfs/model"
	"github.com/urfave/cli"
)

var (
	VERSION string
)

func main() {
	app := cli.NewApp()
	app.Name = "rdbfs"
	app.Usage = "FUSE filesystem implemented with relational databases"
	app.Version = VERSION
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "db-driver",
			Usage:  "set database driver type",
			Value:  "mysql",
			EnvVar: "DB_DRIVER",
		},
		cli.StringFlag{
			Name:   "db-url",
			Usage:  "set database connection url",
			EnvVar: "DB_URL",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:      "start",
			Usage:     "start the fuse server",
			ArgsUsage: "mount-path",
			Action:    startAction,
		},
		{
			Name:   "init",
			Usage:  "Initialize db table",
			Action: initAction,
		},
		{
			Name:   "drop",
			Usage:  "Drop db table",
			Action: dropAction,
		},
	}
	app.Run(os.Args)
}

func getDao(c *cli.Context) (*model.Dao, error) {
	dbDriver := c.GlobalString("db-driver")
	if dbDriver == "" {
		return nil, errors.New("no db driver given, see rdbfs start --help")
	}

	dbUrl := c.GlobalString("db-url")
	if dbUrl == "" {
		return nil, errors.New("no db url given, see rdbfs start --help")
	}

	dao, err := model.NewDao(dbDriver, dbUrl)
	if err != nil {
		return nil, err
	}

	return dao, nil
}

func startAction(c *cli.Context) error {
	mountPath := c.Args().Get(0)
	if mountPath == "" {
		return cli.NewExitError("no mount path given, see rdbfs start --help", 1)
	}

	dao, err := getDao(c)
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	defer dao.Close()

	if err := dao.AutoMigrate(); err != nil {
		return cli.NewExitError(err, 1)
	}

	nfs := pathfs.NewPathNodeFs(&filesystem.RdbFs{FileSystem: pathfs.NewDefaultFileSystem(), Dao: dao}, nil)
	server, _, err := nodefs.MountRoot(mountPath, nfs.Root(), nil)
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	server.Serve()

	return nil
}

func initAction(c *cli.Context) error {
	dao, err := getDao(c)
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	defer dao.Close()

	if err := dao.AutoMigrate(); err != nil {
		return cli.NewExitError(err, 1)
	}

	return nil
}

func dropAction(c *cli.Context) error {
	dao, err := getDao(c)
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	defer dao.Close()

	if err := dao.DropTable(); err != nil {
		return cli.NewExitError(err, 1)
	}

	return nil
}

module github.com/ghetzel/go-stockutil

go 1.17

require (
	github.com/gobwas/glob v0.2.3
	k8s.io/apimachinery v0.23.5
	k8s.io/client-go v11.0.0+incompatible
)

replace (
	k8s.io/api => k8s.io/api v0.19.1
	k8s.io/apimachinery => k8s.io/apimachinery v0.19.1
	k8s.io/client-go => k8s.io/client-go v0.19.1
)

require (
	github.com/alecthomas/assert v0.0.0-20170929043011-405dbfeb8e38
	github.com/alecthomas/colour v0.1.0 // indirect
	github.com/alecthomas/repr v0.0.0-20201120212035-bb82daffcca2 // indirect
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/davecgh/go-spew v1.1.1
	github.com/dsnet/compress v0.0.1
	github.com/erikstmartin/go-testdb v0.0.0-20160219214506-8d10e4a1bae5 // indirect
	github.com/fatih/structs v1.1.0
	github.com/ghetzel/go-defaults v1.2.0
	github.com/ghetzel/testify v1.4.1
	github.com/ghetzel/uuid v0.0.0-20171129191014-dec09d789f3d
	github.com/grandcat/zeroconf v1.0.0
	github.com/h2non/filetype v1.1.3
	github.com/hashicorp/go-multierror v1.1.1
	github.com/jackpal/gateway v1.0.7
	github.com/jbenet/go-base58 v0.0.0-20150317085156-6237cf65f3a6
	github.com/jdkato/prose v1.2.1
	github.com/jdxcode/netrc v0.0.0-20210204082910-926c7f70242a
	github.com/jlaffaye/ftp v0.0.0-20220310202011-d2c44e311e78
	github.com/juliangruber/go-intersect v1.1.0
	github.com/kellydunn/golang-geo v0.7.0
	github.com/kylelemons/go-gypsy v1.0.0 // indirect
	github.com/lib/pq v1.10.4 // indirect
	github.com/martinlindhe/unit v0.0.0-20210313160520-19b60e03648d
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14
	github.com/mattn/go-shellwords v1.0.12
	github.com/melbahja/goph v1.3.0
	github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d
	github.com/mitchellh/mapstructure v1.4.3
	github.com/op/go-logging v0.0.0-20160315200505-970db520ece7
	github.com/phayes/freeport v0.0.0-20220201140144-74d24b5ae9f5
	github.com/pkg/sftp v1.13.4
	github.com/sergi/go-diff v1.1.0 // indirect
	github.com/urfave/negroni v1.0.0
	github.com/ziutek/mymysql v1.5.4 // indirect
	golang.org/x/crypto v0.0.0-20220331220935-ae2d96664a29
	golang.org/x/net v0.0.0-20220403103023-749bd193bc2b
	gopkg.in/neurosnap/sentences.v1 v1.0.7 // indirect
)

require (
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/kr/fs v0.1.0 // indirect
	github.com/miekg/dns v1.1.48 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/mod v0.6.0-dev.0.20220106191415-9b9b3d81d5e3 // indirect
	golang.org/x/sys v0.0.0-20220403205710-6acee93ad0eb // indirect
	golang.org/x/term v0.0.0-20210927222741-03fcf44c2211 // indirect
	golang.org/x/tools v0.1.10 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	gopkg.in/yaml.v2 v2.2.8 // indirect
)

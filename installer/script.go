package installer

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
)

var InstallScript = `#!/bin/bash -x
groupadd {{appname}}
useradd {{appname}} -g {{appname}} -M -s /sbin/nologin
install -m 777 ./{{appname}} /usr/local/bin/{{appname}} 
test -d /usr/lib/systemd/system || mkdir -p /usr/lib/systemd/system
cat>/usr/lib/systemd/system/{{appname}}.service<<EOF
[Unit]
Description={{appname}}
After=network.target

[Service]
LimitNOFILE=65535
LimitNPROC=65535
User={{appname}}
ExecStart=/usr/local/bin/{{appname}}

[Install]
WantedBy=multi-user.target
EOF
chmod 600 /usr/lib/systemd/system/{{appname}}.service
systemctl enable {{appname}} && systemctl daemon-reload
`

func Install(appname string) error {
	InstallScript = strings.ReplaceAll(InstallScript, "{{appname}}", appname)
	scriptfile := fmt.Sprintf("/tmp/%s_install.sh", appname)
	_ = ioutil.WriteFile(scriptfile, []byte(InstallScript), 777)
	if err := exec.Command("/bin/bash", scriptfile).Run(); err != nil {
		return err
	}
	return nil
}

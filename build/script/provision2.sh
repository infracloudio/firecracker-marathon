set -ex

cat <<EOF > /root/server.go
package main 

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"plugin"

	"./context"
)

const (
	CODE_PATH = "/userfunc/user"
)

type (
	FunctionLoadRequest struct {
		// FilePath is an absolute filesystem path to the
		// function. What exactly is stored here is
		// env-specific. Optional.
		FilePath string `json:"filepath"`

		// FunctionName has an environment-specific meaning;
		// usually, it defines a function within a module
		// containing multiple functions. Optional; default is
		// environment-specific.
		FunctionName string `json:"functionName"`

		// URL to expose this function at. Optional; defaults
		// to "/".
		URL string `json:"url"`
	}
)

var userFunc http.HandlerFunc

func loadPlugin(codePath, entrypoint string) http.HandlerFunc {

	// if codepath's a directory, load the file inside it
	info, err := os.Stat(codePath)
	if err != nil {
		panic(err)
	}
	if info.IsDir() {
		files, err := ioutil.ReadDir(codePath)
		if err != nil {
			panic(err)
		}
		if len(files) == 0 {
			panic("No files to load")
		}
		fi := files[0]
		codePath = filepath.Join(codePath, fi.Name())
	}

	fmt.Printf("loading plugin from %v\n", codePath)
	p, err := plugin.Open(codePath)
	if err != nil {
		panic(err)
	}
	sym, err := p.Lookup(entrypoint)
	if err != nil {
		panic("Entry point not found")
	}

	switch h := sym.(type) {
	case *http.Handler:
		return (*h).ServeHTTP
	case *http.HandlerFunc:
		return *h
	case func(http.ResponseWriter, *http.Request):
		return h
	case func(context.Context, http.ResponseWriter, *http.Request):
		return func(w http.ResponseWriter, r *http.Request) {
			c := context.New()
			h(c, w, r)
		}
	default:
		panic("Entry point not found: bad type")
	}
}

func specializeHandler(w http.ResponseWriter, r *http.Request) {
	if userFunc != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Not a generic container"))
		return
	}

	_, err := os.Stat(CODE_PATH)
	if err != nil {
		if os.IsNotExist(err) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(CODE_PATH + ": not found"))
			return
		} else {
			panic(err)
		}
	}

	fmt.Println("Specializing ...")
	userFunc = loadPlugin(CODE_PATH, "Handler")
	fmt.Println("Done")
}

func specializeHandlerV2(w http.ResponseWriter, r *http.Request) {
	if userFunc != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Not a generic container"))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var loadreq FunctionLoadRequest
	err = json.Unmarshal(body, &loadreq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = os.Stat(loadreq.FilePath)
	if err != nil {
		if os.IsNotExist(err) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(CODE_PATH + ": not found"))
			return
		} else {
			panic(err)
		}
	}

	fmt.Println("Specializing ...")
	userFunc = loadPlugin(loadreq.FilePath, loadreq.FunctionName)
	fmt.Println("Done")
}

func readinessProbeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func main() {
	http.HandleFunc("/healthz", readinessProbeHandler)
	http.HandleFunc("/specialize", specializeHandler)
	http.HandleFunc("/v2/specialize", specializeHandlerV2)

	// Generic route -- all http requests go to the user function.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if userFunc == nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Generic container: no requests supported"))
			return
		}
		userFunc(w, r)
	})

	fmt.Println("Listening on 8888 ...")
	http.ListenAndServe(":8888", nil)
}

EOF

mkdir /root/context
cat <<EOF > /root/context/context.go
package context

type (
	Context map[string]interface{}
)

func New() Context {
	ctx := make(map[string]interface{})
	return ctx
}
EOF



cat <<EOF > /etc/init.d/go-server 
#!/bin/sh
### BEGIN INIT INFO
# Provides: go-server
# Required-Start:
# Required-Stop:
# Default-Start: 2 3 4 5
# Default-Stop: 0 1 6
# Short-Description: start go-server at boot time
### END INIT INFO

PATH=/usr/local/sbin:/usr/local/bin:/sbin:/bin:/usr/sbin:/usr/bin:
DAEMON=go-server
NAME="go-server"
DESC="go-server"


. /lib/lsb/init-functions


set -e

# Carry out specific functions when asked to by the system
case "\$1" in
  start)
    echo "Starting script blah "
    echo "Could do more here"
    go-server &
    #start-stop-daemon --start --background --quiet --exec $DAEMON
    ;;
  stop)
    echo "Stopping script blah"
    echo "Could do more here"
    ;;
  *)
    echo "Usage: /etc/init.d/blah {start|stop}"
    exit 1
    ;;
esac

exit 0
EOF


cat <<EOF > /etc/init.d/infra-start
#!/bin/sh
### BEGIN INIT INFO
# Provides: infra-start
# Required-Start:
# Required-Stop:
# Default-Start: 2 3 4 5
# Default-Stop: 0 1 6
# Short-Description: start infra-start at boot time
### END INIT INFO

PATH=/usr/local/sbin:/usr/local/bin:/sbin:/bin:/usr/sbin:/usr/bin:



. /lib/lsb/init-functions


set -e

# Carry out specific functions when asked to by the system
case "\$1" in
  start)
    echo "Starting script blah "
    echo "Could do more here"
		mount /dev/vdb /mnt
		line=\$(head -n 1 /mnt/ip.txt)
		ip addr flush dev eth0
		ip link set dev eth0 up
		ip addr add dev eth0 \$line
		ip route add default via 172.17.0.1
    ;;
  stop)
    echo "Stopping script blah"
    echo "Could do more here"
    ;;
  *)
    echo "Usage: /etc/init.d/blah {start|stop}"
    exit 1
    ;;
esac

exit 0
EOF

/usr/local/go/bin/go build -o /root/go-server /root/server.go
#chmod +x /root/go-server
mv /root/go-server /usr/local/bin
chmod 755 /etc/init.d/go-server
chmod 755 /etc/init.d/infra-start
update-rc.d infra-start defaults 20 03
update-rc.d go-server defaults 97 03
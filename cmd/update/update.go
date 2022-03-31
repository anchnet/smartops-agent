package update

import (
	"bytes"
	"crypto"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	http1 "github.com/anchnet/smartops-agent/pkg/http"
	"github.com/anchnet/smartops-agent/pkg/util/conf"
	log "github.com/cihub/seelog"
	"github.com/kardianos/osext"
)

func execPath() (string, error) {
	here, err := osext.Executable()
	if err != nil {
		return "", err
	}
	return filepath.EvalSymlinks(here)
}

type Options struct {
	Sha256Sum string
	Mode      os.FileMode //0755
}

func Update(remoteAgent io.Reader, opts Options) error {
	if opts.Mode == 0 {
		opts.Mode = 0777
	}
	var err error
	localPath, err := execPath()
	if err != nil {
		return err
	}

	var newBytes []byte
	if newBytes, err = ioutil.ReadAll(remoteAgent); err != nil {
		return err
	}

	if err = checksum(newBytes, opts.Sha256Sum); err != nil {
		return err
	}
	localDir := filepath.Dir(localPath)
	filename := filepath.Base(localPath)

	newPath := filepath.Join(localDir, fmt.Sprintf(".%s.new", filename))
	fp, err := os.OpenFile(newPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, opts.Mode)
	if err != nil {
		return err
	}
	defer fp.Close()

	_, err = io.Copy(fp, bytes.NewReader(newBytes))
	if err != nil {
		return err
	}
	fp.Close()

	oldPath := filepath.Join(localDir, fmt.Sprintf(".%s.old", filename))
	_ = os.Remove(oldPath)

	err = os.Rename(localPath, oldPath)
	if err != nil {
		return err
	}
	err = os.Rename(newPath, localPath)
	if err != nil {
		rerr := os.Rename(oldPath, localPath)
		if rerr != nil {
			return errors.New("failed, waiting for the next update")
		}
		return err
	}

	return nil
}

func checksum(agent []byte, remoteSum string) error {
	sha := crypto.SHA256
	hash := sha.New()
	_, err := hash.Write(agent)
	if err != nil {
		return err
	}
	sum := hash.Sum([]byte{})
	sumS := fmt.Sprintf("%x", sum)
	if !strings.EqualFold(sumS, remoteSum) {
		return fmt.Errorf("agent checksum wrong. local calc: %x, get: %x", sum, remoteSum)
	}
	return nil
}

func GetUpdateInfo(data interface{}) error {
	serverinfo := http1.GetServerInfoData()

	if serverinfo.AutoUpdate.URL == "" {
		return errors.New("unable to upgrade currently")
	}

	url := serverinfo.AutoUpdate.URL
	req, err := http.NewRequest(serverinfo.AutoUpdate.Method, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return errors.New("call get metric status not 200")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("ioutil read Error")
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return errors.New("json Unmarshal Error")
	}
	return nil
}

func AutoUpdate() error {
	type Data struct {
		Data struct {
			Addr      string `json:"addr"`
			Sha256sum string `json:"sha256sum"`
			Domain    string `json:"domain"`
		} `json:"data"`
	}
	data := Data{}

	for {
		time.Sleep(10 * time.Second)
		log.Info("start try update agent")
		err := GetUpdateInfo(&data)
		if err != nil {
			log.Error("get update info err: ", err)
			continue
		}

		if data.Data.Addr == "" || data.Data.Sha256sum == "" {
			log.Info("no update required")
			continue
		}

		resp, err := http.Get(data.Data.Addr)
		if err != nil {
			log.Error(err)
			continue
		}

		if resp.StatusCode != 200 {
			log.Error("get agent binary status code not 200")
			continue
		}

		err = Update(resp.Body, Options{
			Sha256Sum: data.Data.Sha256sum,
		})

		resp.Body.Close()
		if err != nil {
			log.Error(err)
			continue
		}
		//restrat agent by daemon
		log.Info("agent update success")
		if data.Data.Domain != "" {
			if err := conf.ChangeConfSite(data.Data.Domain, "./conf/smartops.yaml"); err != nil {
				log.Error("change smartops site error", err)
			}
		}
		time.Sleep(1 * time.Second)
		os.Exit(-1)
	}
}

func OnceUpdate() error {
	type Data struct {
		Data struct {
			Addr      string `json:"addr"`
			Sha256sum string `json:"sha256sum"`
			Domain    string `json:"domain"`
		} `json:"data"`
	}
	data := Data{}

	log.Info("start try update agent")
	err := GetUpdateInfo(&data)
	if err != nil {
		log.Error("get update info err: ", err)
		return err
	}

	if data.Data.Addr == "" || data.Data.Sha256sum == "" {
		log.Info("no update required")
		return err
	}

	resp, err := http.Get(data.Data.Addr)
	if err != nil {
		log.Error(err)
		return err
	}

	if resp.StatusCode != 200 {
		log.Error("get agent binary status code not 200")
		return err
	}

	err = Update(resp.Body, Options{
		Sha256Sum: data.Data.Sha256sum,
	})

	resp.Body.Close()
	if err != nil {
		log.Error(err)
		return err
	}
	//restrat agent by daemon
	log.Info("agent update success")
	if data.Data.Domain != "" {
		if err := conf.ChangeConfSite(data.Data.Domain, "./conf/smartops.yaml"); err != nil {
			log.Error("change smartops site error", err)
		}
	}
	time.Sleep(1 * time.Second)
	os.Exit(-1)
	return nil
}

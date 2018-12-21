package helmcontroller

import (
	"io/ioutil"
	"net/url"

	"k8s.io/helm/pkg/getter"
	helm_env "k8s.io/helm/pkg/helm/environment"
)

var settings helm_env.EnvSettings

// readFile load a file from the local directory or a remote file with a url
func readFile(filePath string) ([]byte, error) {
	u, _ := url.Parse(filePath)
	p := getter.All(settings)

	getterConstructor, err := p.ByScheme(u.Scheme)

	if err != nil {
		return ioutil.ReadFile(filePath)
	}

	getter, err := getterConstructor(filePath, "", "", "")
	if err != nil {
		return []byte{}, err
	}
	data, err := getter.Get(filePath)
	return data.Bytes(), err
}

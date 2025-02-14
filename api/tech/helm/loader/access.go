package loader

import (
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"

	"ocm.software/ocm/api/tech/helm"
	"ocm.software/ocm/api/utils/blobaccess/blobaccess"
)

type accessLoader struct {
	access helm.ChartAccess
}

func AccessLoader(acc helm.ChartAccess) Loader {
	return &accessLoader{access: acc}
}

func (l *accessLoader) Close() error {
	return l.access.Close()
}

func (l *accessLoader) ChartArchive() (blobaccess.BlobAccess, error) {
	return l.access.Chart()
}

func (l *accessLoader) ChartArtefactSet() (blobaccess.BlobAccess, error) {
	return l.access.ArtefactSet()
}

func (l *accessLoader) Chart() (*chart.Chart, error) {
	acc, err := l.access.Chart()
	if err != nil {
		return nil, err
	}
	defer acc.Close()
	r, err := acc.Reader()
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return loader.LoadArchive(r)
}

func (l *accessLoader) Provenance() ([]byte, error) {
	prov, err := l.access.Prov()
	if prov == nil || err != nil {
		return nil, err
	}
	defer prov.Close()
	return prov.Get()
}

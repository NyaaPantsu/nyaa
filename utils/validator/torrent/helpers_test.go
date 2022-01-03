package torrentValidator

import (
	"testing"

	"github.com/anacrolix/torrent/metainfo"
	"github.com/stretchr/testify/assert"
)

func TestCheckTrackers(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		Info         metainfo.MetaInfo
		Expected     []string
		ExpectedList metainfo.AnnounceList
	}{
		{
			metainfo.MetaInfo{
				AnnounceList: [][]string{{"udp://nn.fof:4545/ano", "lol"}, {"://nyaa.tr/ann"}, {".co", "http://mont.co:444/ann"}},
			},
			[]string{"udp://nn.fof:4545/ano", "http://mont.co:444/ann", "udp://tracker.uw0.xyz:6969/announce", "http://anidex.moe:6969/announce"},
			[][]string{{"udp://nn.fof:4545/ano"}, {"http://mont.co:444/ann"}, {"udp://tracker.uw0.xyz:6969/announce"}, {"http://anidex.moe:6969/announce"}},
		},
		{
			metainfo.MetaInfo{
				AnnounceList: [][]string{{"http://open.nyaatorrents.info:6544", "http://mont.co:444/ann"}}, // dead tracker
			},
			[]string{"http://mont.co:444/ann", "udp://tracker.uw0.xyz:6969/announce", "http://anidex.moe:6969/announce"},
			[][]string{{"http://mont.co:444/ann"}, {"udp://tracker.uw0.xyz:6969/announce"}, {"http://anidex.moe:6969/announce"}},
		},
		{
			metainfo.MetaInfo{
				AnnounceList: [][]string{{"http://open.nyaatorrents.info:6544", "http://mont.co:444/ann", "http://mont.co:4434/ann"}, {"http://mont.co:444/anno"}}, // dead tracker
			},
			[]string{"http://mont.co:444/ann", "http://mont.co:4434/ann", "http://mont.co:444/anno", "udp://tracker.uw0.xyz:6969/announce", "http://anidex.moe:6969/announce"},
			[][]string{{"http://mont.co:444/ann", "http://mont.co:4434/ann"}, {"http://mont.co:444/anno"}, {"udp://tracker.uw0.xyz:6969/announce"}, {"http://anidex.moe:6969/announce"}},
		},
	}

	for _, test := range tests {
		trackers := CheckTrackers(&test.Info)
		assert.Equal(test.Expected, trackers)
		assert.Equal(test.ExpectedList, test.Info.AnnounceList)
	}

}

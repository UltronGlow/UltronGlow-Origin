package alien

import (
	"github.com/UltronGlow/UltronGlow-Origin/common"
	"github.com/UltronGlow/UltronGlow-Origin/log"
	"strconv"
	"strings"
)

func verifyPocStringV1(block, nonce, blockhash, pocstr, roothash, deviceAddr string) bool {
	poc := strings.Split(pocstr, ",")
	//log.Info("verifyPocStringV1","pocstr",pocstr)
	if len(poc) < 10 {
		log.Warn("verifyStoragePoc", "invalide poc string format")
		return false
	}
	if poc[0] != "v1" {
		log.Warn("verifyStoragePocV1", "invalide version tag")
		return false
	}

	sampleNumberpos := 4
	blocknumberpos := 6
	b0pos := 8
	b1pos := 9
	bnpos := 10

	if !verifyB0(block, nonce, blockhash, poc[b0pos],deviceAddr ) {
		log.Warn("verifyPocString", "verify b0 failed",pocstr)
		return false
	}

	if !verifyBn(poc[sampleNumberpos], poc[b0pos], poc[bnpos], poc[b1pos]) {
		log.Warn("verifyPocString", "verify bn failed")
		return false
	}

	n, _ := strconv.ParseUint(poc[sampleNumberpos], 10, 64)
	if !verifySamplePos(n, poc[1], poc[2], poc[3], poc[blocknumberpos]) {
		log.Warn("verifyPocString", "verify samplenumber failed")
		return false
	}

	if n&1 != 0 {
		h1 := Hash(poc[b1pos], poc[bnpos], "")
		return verifyPocV1(poc[11:], h1, roothash, n)
	} else {
		h1 := Hash(poc[bnpos], poc[bnpos+1], "")
		return verifyPocV1(poc[12:], h1, roothash, n)
	}
}

func verifyStoragePocV1(pocstr, roothash string, nonce uint64) bool {
	poc := strings.Split(pocstr, ",")
	if len(poc) < 10 {
		log.Warn("verifyStoragePocV1", "invalide v1 poc string format")
		return false
	}
	if poc[0] != "v1" {
		log.Warn("verifyStoragePocV1", "invalide version tag")
		return false
	}
	if poc[2] != strconv.FormatUint(nonce, 10) {
		log.Warn("verifyStoragePocV1", "invalide nonce")
		return false
	}

	sampleNumberpos := 4
	blocknumberpos := 6
	b0pos := 8
	b1pos := 9
	bnpos := 10

	if !verifyBn(poc[sampleNumberpos], poc[b0pos], poc[bnpos], poc[b1pos]) {
		log.Warn("verifyPocStringV1", "verify bn failed")
		return false
	}

	n, _ := strconv.ParseUint(poc[sampleNumberpos], 10, 64)
	if !verifySamplePos(n, poc[1], poc[2], poc[3], poc[blocknumberpos]) {
		log.Warn("verifyPocStringV1", "verify samplenumber failed")
		return false
	}

	if n&1 != 0 {
		h1 := Hash(poc[b1pos], poc[bnpos], "")
		return verifyPocV1(poc[11:], h1, roothash, n)
	} else {
		h1 := Hash(poc[bnpos], poc[bnpos+1], "")
		return verifyPocV1(poc[12:], h1, roothash, n)
	}
}

func verifyPocV1(pocstr []string, h1, roothash string, r uint64) bool {
	var (
		hash  string
		round int
		hashpos int
	)

	r = r / 2
	hashpos = int(r & 1)
	hash = h1
	for i := 0; i < len(pocstr); i++ {
		if i+1 >= len(pocstr) {
			break
		}
		if round&1 != 1 {
			if hashpos == 0 {
				hash = Acc(hash, pocstr[i], "")
			} else {
				hash = Acc(pocstr[i], hash, "")
			}
		} else {
			if hashpos == 0 {
				hash = Hash(hash, pocstr[i], "")
			} else {
				hash = Hash(pocstr[i], hash, "")
			}
		}

		r = r / 2
		hashpos = int(r & 1)
		round++
	}
	if hash == pocstr[len(pocstr)-1] && common.HexToHash(hash) == common.HexToHash(roothash) {
		return true
	}
	log.Warn("verifyPocV1", "root hash:", hash, "roothash", roothash, "pocstr[len(pocstr)-1]", pocstr[len(pocstr)-1], "common.HexToHash(hash)", common.HexToHash(hash))
	return false
}

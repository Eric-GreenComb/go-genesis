// Copyright 2016 The go-daylight Authors
// This file is part of the go-daylight library.
//
// The go-daylight library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-daylight library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-daylight library. If not, see <http://www.gnu.org/licenses/>.

package parser

import (
	"fmt"
	"strings"

	"github.com/DayLightProject/go-daylight/packages/consts"
	"github.com/DayLightProject/go-daylight/packages/script"
	"github.com/DayLightProject/go-daylight/packages/smart"
)

func init() {
	smart.Extend(&script.ExtendData{map[string]interface{}{
		"DBInsert": DBInsert,
	}, map[string]string{
		`*parser.Parser`: `parser`,
	}})
}

func (p *Parser) getExtend() *map[string]interface{} {
	head := p.TxPtr.(*consts.TXHeader) //consts.HeaderNew(contract.parser.TxPtr)
	var citizenId, walletId int64
	if head.StateId > 0 {
		citizenId = head.UserId
	} else {
		walletId = head.UserId
	}
	// test
	if citizenId == 1 {
		walletId = 1
	}
	block := int64(0)
	if p.BlockData != nil {
		block = p.BlockData.BlockId
	}
	extend := map[string]interface{}{`type`: head.Type, `time`: head.Type, `state`: head.StateId,
		`block`: block, `citizen`: citizenId, `wallet`: walletId,
		`parser`: p, `contract`: p.TxContract}
	for key, val := range p.TxData {
		extend[key] = val
	}
	/*	v := reflect.ValueOf(contract.parser.TxPtr).Elem()
		t := v.Type()
		for i := 1; i < t.NumField(); i++ {
			extend[t.Field(i).Name] = v.Field(i).Interface()
		}*/
	//fmt.Println(`Extend`, extend)
	return &extend
}

func (p *Parser) CallContract(flags int) (err error) {
	methods := []string{`init`, `front`, `main`}
	extend := p.getExtend()
	for i := uint32(0); i < 3; i++ {
		if (flags & (1 << i)) > 0 {
			cfunc := p.TxContract.GetFunc(methods[i])
			if cfunc == nil {
				continue
			}
			p.TxContract.Called = 1 << i
			_, err = smart.Run(cfunc, nil, extend)
			if err != nil {
				return
			}
		}
	}
	return
}

func DBInsert(p *Parser, tblname string, params string, val ...interface{}) (err error) { // map[string]interface{}) {
	fmt.Println(`DBInsert`, tblname, params, val, len(val))
	err = p.selectiveLoggingAndUpd(strings.Split(params, `,`), val, tblname, nil, nil, true)
	return
}
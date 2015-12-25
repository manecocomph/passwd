package impt

import (
	"strings"
)

const SEP_WHITE_SPACE = "WHITE_SPACE"

type parseFunc func(string, *ImportMetaData) Any

type ImportMetaData struct {
	Name        string
	FileName    string //Can be file name or folder name
	Sep         string
	EmailSeq    int
	PasswdSeq   int
	UserNameSeq int //negative number: not exist
}

func mapParseFunc(metaData *ImportMetaData) func(string, *ImportMetaData) Any {
	if SEP_WHITE_SPACE == metaData.Sep {
		if 0 > metaData.UserNameSeq {
			return convertEmailPwdByWhiteSpace
		} else {
			return convertEmailPwdUserNameByWhiteSpace
		}
	} else {
		if 0 > metaData.UserNameSeq {
			return convertEmailPwdBySep
		} else {
			return convertEmailPwdUserNameBySep
		}
	}
}

func convertEmailPwdByWhiteSpace(line string, metaData *ImportMetaData) Any {
	splits := strings.Fields(line)
	if 2 == len(splits) {
		return Account{line, splits[metaData.EmailSeq], splits[metaData.PasswdSeq], ""}
	} else {
		return Account{Raw: line}
	}
}

func convertEmailPwdBySep(line string, metaData *ImportMetaData) Any {
	splits := strings.Split(line, metaData.Sep)
	if 2 == len(splits) {
		return Account{line, splits[metaData.EmailSeq], splits[metaData.PasswdSeq], ""}
	} else {
		return Account{Raw: line}
	}
}

func convertEmailPwdUserNameByWhiteSpace(line string, metaData *ImportMetaData) Any {
	splits := strings.Fields(line)
	if 3 == len(splits) {
		return Account{line, splits[metaData.EmailSeq], splits[metaData.PasswdSeq], splits[metaData.UserNameSeq]}
	} else {
		return Account{Raw: line}
	}
}

func convertEmailPwdUserNameBySep(line string, metaData *ImportMetaData) Any {
	splits := strings.Split(line, metaData.Sep)
	if 3 == len(splits) {
		return Account{line, splits[metaData.EmailSeq], splits[metaData.PasswdSeq], splits[metaData.UserNameSeq]}
	} else {
		return Account{Raw: line}
	}
}

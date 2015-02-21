package moderniejapanizer

import (
	"strconv"
	"syscall"

	"github.com/hnakamur/w32registry"
	"github.com/hnakamur/w32syscall"
)

const JapaneseLanguagePackUpdateID = "00a156d4-3876-4cd5-bd38-517679c6ba59"

// https://msdn.microsoft.com/en-us/library/dd374073.aspx
const JapaneseLocationCode = "122"

// http://social.technet.microsoft.com/wiki/contents/articles/6281.how-to-set-the-keyboard-layout-through-group-policy-gpo.aspx
const (
	EnglishUnitedStatesKeyboardCode = "00000409"
	JapaneseJapanKeyboardCode       = "00000411"
)

const JapaneseDisplayLanguageCode = "ja-JP"

// https://technet.microsoft.com/en-us/library/cc978632.aspx
// http://support.microsoft.com/kb/102978/
var JapaneseLanguageAndRegionalFormats = map[string]string{
	"iCalendarType":    "1",
	"iCountry":         "81",
	"iCurrDigits":      "0",
	"iCurrency":        "0",
	"iDate":            "2",
	"iDigits":          "2",
	"iFirstDayOfWeek":  "6",
	"iFirstWeekOfYear": "0",
	"iLZero":           "1",
	"iMeasure":         "0",
	"iNegCurr":         "1",
	"iNegNumber":       "1",
	"iPaperSize":       "9",
	"iTime":            "1",
	"iTimePrefix":      "0",
	"iTLZero":          "0",
	"Locale":           "00000411",
	"LocaleName":       "ja-JP",
	"NumShape":         "1",
	"s1159":            "午前",
	"s2359":            "午後",
	"sCountry":         "Japan",
	"sCurrency":        "\xC2\xA5",
	"sDate":            "/",
	"sDecimal":         ".",
	"sGrouping":        "3;0",
	"sLanguage":        "JPN",
	"sList":            ",",
	"sLongDate":        "yyyy'年'M'月'd'日'",
	"sMonDecimalSep":   ".",
	"sMonGrouping":     "3;0",
	"sMonThousandSep":  ",",
	"sNativeDigits":    "0123456789",
	"sNegativeSign":    "-",
	"sPositiveSign":    "",
	"sShortDate":       "yyyy/MM/dd",
	"sThousand":        ",",
	"sTime":            ":",
	"sTimeFormat":      "H:mm:ss",
	"sShortTime":       "H:mm",
	"sYearMonth":       "yyyy'年'M'月'",
}

func SetLocation(locationCode string) error {
	return w32registry.SetKeyValueString(syscall.HKEY_CURRENT_USER, `Control Panel\International\Geo`, "Nation", locationCode)
}

func SetKeyboards(keyboardCodes []string) error {
	for i, code := range keyboardCodes {
		err := w32registry.SetKeyValueString(syscall.HKEY_CURRENT_USER, `Keyboard Layout\Preload`, strconv.Itoa(i+1), code)
		if err != nil {
			return err
		}
	}
	return nil
}

func SetLanguageAndRegionalFormats(formats map[string]string) error {
	for k, v := range formats {
		err := w32registry.SetKeyValueString(syscall.HKEY_CURRENT_USER, `Control Panel\International`, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func Reboot(reason uint32) error {
	process, err := syscall.GetCurrentProcess()
	if err != nil {
		return err
	}

	var hToken syscall.Token
	err = syscall.OpenProcessToken(process, syscall.TOKEN_ADJUST_PRIVILEGES|syscall.TOKEN_QUERY, &hToken)
	if err != nil {
		return err
	}

	seShutdownNameP, err := syscall.UTF16PtrFromString(w32syscall.SE_SHUTDOWN_NAME)
	if err != nil {
		return err
	}
	var tkp w32syscall.TokenPrivileges
	err = w32syscall.LookupPrivilegeValue(nil, seShutdownNameP, &tkp.Privileges[0].Luid)
	if err != nil {
		return err
	}
	tkp.PrivilegeCount = 1
	tkp.Privileges[0].Attributes = w32syscall.SE_PRIVILEGE_ENABLED

	err = w32syscall.AdjustTokenPrivileges(hToken, false, &tkp, 0, nil, nil)
	if err != nil {
		return err
	}

	err = w32syscall.ExitWindowsEx(w32syscall.EWX_REBOOT, reason)
	if err != nil {
		return err
	}

	tkp.Privileges[0].Attributes = 0
	return w32syscall.AdjustTokenPrivileges(hToken, false, &tkp, 0, nil, nil)
}

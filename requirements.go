package main

import (
	"log"
	"os"
	"regexp"
	"runtime"
	"strings"
)

func installRequirements() error {
	log.Println("install uiautomator apk")
	if err := installUiautomatorAPK(); err != nil {
		return err
	}
	return installMinicap()
}

func installUiautomatorAPK() error {
	if runtime.GOOS == "windows" {
		return nil
	}
	if checkUiautomatorInstalled() {
		return nil
	}
	baseURL := "https://github.com/openatx/android-uiautomator-server/releases/download/" + apkVersionName
	if _, err := httpDownload("/data/local/tmp/app-debug.apk", baseURL+"/app-uiautomator.apk", 0644); err != nil {
		return err
	}
	if _, err := httpDownload("/data/local/tmp/app-debug-test.apk", baseURL+"/app-uiautomator-test.apk", 0644); err != nil {
		return err
	}
	if err := forceInstallAPK("/data/local/tmp/app-debug.apk"); err != nil {
		return err
	}
	if err := forceInstallAPK("/data/local/tmp/app-debug-test.apk"); err != nil {
		return err
	}
	return nil
}

func installMinicap() error {
	if runtime.GOOS == "windows" {
		return nil
	}
	log.Println("install minicap")
	// if fileExists("/data/local/tmp/minicap") && fileExists("/data/local/tmp/minicap.so") {
	// 	if err := Screenshot("/dev/null"); err != nil {
	// 		log.Println("err:", err)
	// 	} else {
	// 		return nil
	// 	}
	// }
	// remove first to prevent "text file busy"
	os.Remove("/data/local/tmp/minicap")
	os.Remove("/data/local/tmp/minicap.so")

	minicapSource := "https://github.com/codeskyblue/stf-binaries/raw/master/node_modules/minicap-prebuilt/prebuilt"
	propOutput, err := runShell("getprop")
	if err != nil {
		return err
	}
	re := regexp.MustCompile(`\[(.*?)\]:\s*\[(.*?)\]`)
	matches := re.FindAllStringSubmatch(string(propOutput), -1)
	props := make(map[string]string)
	for _, m := range matches {
		var key = m[1]
		var val = m[2]
		props[key] = val
	}
	abi := props["ro.product.cpu.abi"]
	sdk := props["ro.build.version.sdk"]
	pre := props["ro.build.version.preview_sdk"]
	if pre != "" && pre != "0" {
		sdk = sdk + pre
	}
	binURL := strings.Join([]string{minicapSource, abi, "bin", "minicap"}, "/")
	_, err = httpDownload("/data/local/tmp/minicap", binURL, 0755)
	if err != nil {
		return err
	}
	libURL := strings.Join([]string{minicapSource, abi, "lib", "android-" + sdk, "minicap.so"}, "/")
	_, err = httpDownload("/data/local/tmp/minicap.so", libURL, 0644)
	if err != nil {
		return err
	}
	return nil
}

//下载minitouch
func installMinitouch() error {
	baseURL := "https://github.com/codeskyblue/stf-binaries/raw/master/node_modules/minitouch-prebuilt/prebuilt"
	//获取cpu型号
	abi := getCachedProperty("ro.product.cpu.abi")
	//拼接url
	binURL := strings.Join([]string{baseURL, abi, "bin/minitouch"}, "/")
	//通过下载模块将minitouch下载到手机
	_, err := httpDownload("/data/local/tmp/minitouch", binURL, 0755)
	return err
}

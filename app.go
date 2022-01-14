package autoEncode

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type EncodeFactory struct {
	workDir     string
	inDir       string
	outDir      string
	exeDir      string
	tmpDir      string
	profileName string
	Status      EncodeStatuses
	hash        string
	logName     string
	statusName  string
}

func NewEncodeFactory() (EncodeFactory, error) {
	const (
		logName    = "encode.log"
		statusName = "encode_status.csv"
		workDir    = "work"
	)

	factory := EncodeFactory{
		workDir:    workDir,
		logName:    logName,
		statusName: statusName,
		Status:     NewEncodeStatuses(),
		hash:       "",
	}

	if err := factory.refresh(); err != nil {
		return EncodeFactory{}, err
	}

	return factory, nil
}

func checkDir(dir string) error {
	file, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%v:%v", ErrTargetPathNotFound, dir)
		} else {
			return fmt.Errorf("%v:%v", ErrException, dir)
		}
	}
	if !file.IsDir() {
		return fmt.Errorf("%v:%v", ErrTargetIsNotDir, dir)
	}
	return nil
}

func (f *EncodeFactory) Set(inDir, outDir, tmpDir, exeDir, profile string) error {
	if err := checkDir(inDir); err != nil {
		return err
	}
	f.inDir = inDir

	if err := checkDir(outDir); err != nil {
		return err
	}
	f.outDir = outDir

	if err := checkDir(exeDir); err != nil {
		return err
	}
	f.exeDir = exeDir

	if err := checkDir(tmpDir); err != nil {
		return err
	}
	f.tmpDir = tmpDir

	f.profileName = profile
	return nil
}

func (f *EncodeFactory) refresh() error {
	path := filepath.Join(f.workDir, f.statusName)
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%v:%v", ErrFileNotFound, path)
		} else {
			return fmt.Errorf("%v:%v", ErrException, path)
		}
	}

	hash, err := Hash(path)
	if err != nil {
		return ErrHashFailed
	}

	if f.hash != "" && f.hash == hash {
		return nil
	}

	err = f.Status.ReadFile(path)
	if err != nil {
		return fmt.Errorf("%v:%v", ErrReadFileFaield, path)
	}

	f.hash, err = Hash(path)
	if err != nil {
		return err
	}
	return nil
}

func (f *EncodeFactory) update() error {
	path := filepath.Join(f.workDir, f.statusName)
	for {
		err := f.Status.WriteFile(path)
		if err == nil {
			break
		}
		rand.Seed(time.Now().UnixNano())
		n := rand.Intn(90)
		log.Printf("Sleeping %d seconds...\n", n)
		time.Sleep(time.Duration(n) * time.Second)
		fmt.Println("Done")
	}

	hash, err := Hash(path)
	if err != nil {
		return fmt.Errorf("%v:%v", ErrHashFailed, path)
	}

	f.hash = hash
	return nil
}

func (f *EncodeFactory) Add() error {
	log.Println("Add title check start")
	info, err := os.ReadDir(f.inDir)
	if err != nil {
		return fmt.Errorf("%v:%v", ErrReadDirFailed, f.inDir)
	}

	for _, record := range info {
		// ディレクトリは処理対象外
		if record.IsDir() {
			continue
		}

		// "."から始まる隠しファイルは処理対象外
		if strings.HasPrefix(record.Name(), ".") {
			continue
		}

		// 拡張子が.m2ts出ないファイルは処理対象外
		if filepath.Ext(record.Name()) != ".m2ts" {
			continue
		}

		// 対象ファイルがamatsukaze登録済ならばスキップ
		i, _ := f.Status.GetStatus(record.Name())
		if i != Init {
			continue
		}

		// 録画ファイルをローカルディレクトリにコピー
		log.Println("add title:", record.Name())
		src := filepath.Join(f.inDir, record.Name())
		dst := filepath.Join(f.tmpDir, record.Name())
		if err := Copy(src, dst); err != nil {
			return fmt.Errorf("%v:%v", ErrCopyFailed, record.Name())
		}

		// amatsukaze登録処理開始
		path := dst
		if err := f.Status.Add(record.Name()); err != nil {
			log.Printf("[WARN]Add failed:%v ...skiped\n", err)
		}
		if err := f.update(); err != nil {
			return err
		}

		exPath := filepath.ToSlash(filepath.Join(f.exeDir, "AmatsukazeAddTask.exe"))
		_, err := os.Stat(exPath)
		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("%v:%v", ErrFileNotFound, exPath)
			}
			return fmt.Errorf("%v:%v", ErrException, exPath)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		// amatsukaze登録実行
		output, err := exec.CommandContext(ctx, exPath, "-f", path, "-o", f.outDir, "-s", f.profileName).CombinedOutput()
		if err != nil {
			return fmt.Errorf("%v-%v", ErrFailedCMD, err)
		}
		if err := f.Status.Set(record.Name(), Add); err != nil {
			return err
		}
		if err := f.update(); err != nil {
			return err
		}

		// 登録処理結果確認
		str, err := BytesFromShiftJIS(output)
		if err != nil {
			return err
		}
		if !strings.HasSuffix(strings.TrimSpace(str), "1件追加しました") {
			return fmt.Errorf("%v:%v", ErrAddAmatsukaze, str)
		}

		if err := f.Status.Set(record.Name(), Added); err != nil {
			return fmt.Errorf("%v:%v:%v", ErrSetFailed, err, record.Name())
		}
		if err := f.update(); err != nil {
			return fmt.Errorf("%v:%v:%v", ErrUpdateFailed, err, record.Name())
		}
		log.Println("add success:", record.Name())
	}

	return nil
}

func (f *EncodeFactory) Start(title string) error {
	if err := f.refresh(); err != nil {
		return err
	}

	i, err := f.Status.GetStatus(title)
	if err != nil {
		return fmt.Errorf("%v:'%v'", err, title)
	}
	if i != Added {
		return fmt.Errorf("%v:%v", ErrStatusFailed, i)
	}

	if err := f.Status.Set(title, Started); err != nil {
		return fmt.Errorf("%v:%v:%v", ErrSetFailed, err, title)
	}

	if err := f.update(); err != nil {
		return fmt.Errorf("%v:%v:%v", ErrUpdateFailed, err, title)
	}
	return nil
}

func (f *EncodeFactory) Finish(title string) error {
	if err := f.refresh(); err != nil {
		return err
	}

	status, err := f.Status.GetStatus(title)
	if err != nil {
		return fmt.Errorf("%v:%v:'%v'", ErrGetStatusFailed, err, title)
	}
	if status != Started {
		return fmt.Errorf("%v:%v", ErrStatusFailed, status)
	}

	if err := f.Status.Set(title, Finish); err != nil {
		return fmt.Errorf("%v:%v:%v", ErrSetFailed, err, title)
	}

	if err := f.update(); err != nil {
		return fmt.Errorf("%v:%v:%v", ErrUpdateFailed, err, title)
	}
	return nil
}

/**
 * Copyright 2020 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package env

import (
	"io"
	"log"
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

// Config environment variables used by envconfig
type Config struct {
	LogLevel    string `split_words:"true" default:"debug"`
	LogFile     string `split_words:"true"`
	BucketName  string `split_words:"true" required:"true"`
	KmsMkekURI  string `split_words:"true" required:"true"`
	DekPathName string `split_words:"true" required:"true"`
	AAD         string `split_words:"true" required:"true"`
	Client      ClientConfig
	Proxy       ProxyConfig
}

// Logger configures logging based on env variables
func (c Config) Logger() *logrus.Logger {
	badLevelName := false

	level, err := logrus.ParseLevel(c.LogLevel)
	if err != nil {
		badLevelName = true
		level = logrus.DebugLevel
	}

	logger := logrus.New()

	// determine if logs need to be sent to a file in addition to Stdout, which is the default
	if c.LogFile != "" {
		// always overwrite logfile
		logFile, err := os.Create(c.LogFile)
		if err != nil {
			log.Fatal(err)
		}
		mw := io.MultiWriter(os.Stdout, logFile)
		logger.Out = mw
	} else {
		logger.Out = os.Stdout
	}

	logger.Level = level
	if badLevelName {
		logger.WithField("LogLevel", c.LogLevel).WithError(err).Warn("unknown logging level. check your environment variable")
	}

	return logger
}

// Get loads the configuration from environment variables.
func Get() (Config, error) {
	var c Config
	err := envconfig.Process("tinkproxy", &c)
	return c, err
}

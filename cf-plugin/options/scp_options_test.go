package options_test

import (
	"github.com/cloudfoundry-incubator/diego-ssh/cf-plugin/options"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SCPOptions", func() {
	var (
		opts       *options.SCPOptions
		args       []string
		parseError error
	)

	Describe("Parse", func() {
		BeforeEach(func() {
			opts = options.NewSCPOptions()
			args = []string{"scp"}
			parseError = nil
		})

		JustBeforeEach(func() {
			parseError = opts.Parse(args)
		})

		Context("when the command name is missing", func() {
			BeforeEach(func() {
				args = []string{}
			})

			It("returns a usage error", func() {
				Expect(parseError).To(Equal(options.UsageError))
			})
		})

		Context("when the wrong command is specified", func() {
			BeforeEach(func() {
				args = []string{"ssh"}
			})

			It("returns a usage error", func() {
				Expect(parseError).To(Equal(options.UsageError))
			})
		})

		Context("when fewer than two arguments are specified", func() {
			BeforeEach(func() {
				args = append(args, "local.txt")
			})

			It("returns a usage error", func() {
				Expect(parseError).To(MatchError("Source and target must be provided"))
			})
		})

		Context("when two or more arguments are specified", func() {
			BeforeEach(func() {
				args = append(args, "local.txt", "app/1:remote.txt", "app/99:remote.txt")
			})

			It("does not error", func() {
				Expect(parseError).ToNot(HaveOccurred())
			})

			It("sets the first as the source", func() {
				Expect(opts.Sources).To(ConsistOf(
					options.FileLocation{
						Path: "local.txt",
					},
					options.FileLocation{
						AppName: "app",
						Index:   1,
						Path:    "remote.txt",
					},
				))
			})

			It("sets the last argument as the target", func() {
				Expect(opts.Target).To(Equal(options.FileLocation{
					AppName: "app",
					Index:   99,
					Path:    "remote.txt",
				}))
			})
		})

		Context("when the verbose flag is set", func() {
			BeforeEach(func() {
				args = append(args, "-v", "file", "file")
			})

			It("sets the verbose flag", func() {
				Expect(parseError).NotTo(HaveOccurred())
				Expect(opts.Verbose).To(BeTrue())
			})
		})

		Context("when the preserve attributes flag is set", func() {
			BeforeEach(func() {
				args = append(args, "-p", "file", "file")
			})

			It("sets the preserve flag", func() {
				Expect(parseError).NotTo(HaveOccurred())
				Expect(opts.PreserveAttributes).To(BeTrue())
			})
		})

		Context("when the recursive flag is set", func() {
			BeforeEach(func() {
				args = append(args, "-r", "dir", "dir")
			})

			It("sets the recurse flag", func() {
				Expect(parseError).NotTo(HaveOccurred())
				Expect(opts.Recurse).To(BeTrue())
			})
		})
	})

	Describe("ParseLocation", func() {
		Context("when the argument is a local path", func() {
			It("returns a local file location", func() {
				location, err := options.ParseLocation("local.txt")
				Expect(err).NotTo(HaveOccurred())

				Expect(location).To(Equal(options.FileLocation{
					Path: "local.txt",
				}))
			})
		})

		Context("when the argument is a remote path", func() {
			It("returns a remote file location", func() {
				location, err := options.ParseLocation("app:remote.txt")
				Expect(err).NotTo(HaveOccurred())

				Expect(location).To(Equal(options.FileLocation{
					AppName: "app",
					Path:    "remote.txt",
				}))
			})
		})

		Context("when the argument is a remote path with app and index", func() {
			It("returns a remote file location", func() {
				location, err := options.ParseLocation("app/99:remote.txt")
				Expect(err).NotTo(HaveOccurred())

				Expect(location).To(Equal(options.FileLocation{
					AppName: "app",
					Index:   99,
					Path:    "remote.txt",
				}))
			})
		})

		Context("when the argument format is invalid", func() {
			It("returns an error", func() {
				_, err := options.ParseLocation("ap/99/99:remote.txt")
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when the file name contains an escaped colon", func() {
			It("returns an error", func() {
				location, err := options.ParseLocation("file\\:with\\:colons.txt")
				Expect(err).NotTo(HaveOccurred())

				Expect(location).To(Equal(options.FileLocation{
					Path: "file:with:colons.txt",
				}))
			})
		})

		Context("when the location is a windows file", func() {
			It("returns the correct location", func() {
				location, err := options.ParseLocation(`C\:\some\windows\file.txt`)
				Expect(err).NotTo(HaveOccurred())

				Expect(location).To(Equal(options.FileLocation{
					Path: `C:\some\windows\file.txt`,
				}))
			})
		})

		Context("when the location is a UNC name", func() {
			It("returns the correct location", func() {
				location, err := options.ParseLocation(`\\?\D\:\some\windows\file.txt`)
				Expect(err).NotTo(HaveOccurred())

				Expect(location).To(Equal(options.FileLocation{
					Path: `\\?\D:\some\windows\file.txt`,
				}))
			})
		})
	})

	Describe("SCPUsage", func() {
		It("prints usage information", func() {
			usage := options.SCPUsage()

			Expect(usage).To(ContainSubstring("Usage: scp [-prv] [app[/index]:]file1 ... [app[/index]:]file2"))
			Expect(usage).To(ContainSubstring("-p    preserve file times and permissions"))
			Expect(usage).To(ContainSubstring("-r    recurse into directories"))
			Expect(usage).To(ContainSubstring("-v    enable verbose output"))
		})
	})
})

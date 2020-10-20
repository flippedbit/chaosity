package cmd

import (
	"log"
	"os"
	"runtime"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/flippedbit/chaosity/pkg/aws/options"
	"github.com/inconshreveable/go-update"
	"github.com/spf13/cobra"
)

var updateOptions options.UpdateOptions

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Updates chaosity binary to latest version",
	Long:  "In place upgrade for chaosity binary",
	Run: func(cmd *cobra.Command, args []string) {
		sess := session.Must(session.NewSessionWithOptions(session.Options{
			Profile: updateOptions.Profile,
			Config: aws.Config{
				Region: aws.String(updateOptions.Region),
			},
			SharedConfigState: session.SharedConfigEnable,
		}))
		s3svc := s3manager.NewDownloader(sess)

		log.Println("Updating binary... this should only take a few seconds.")

		key := "builds-chaosity/chaosity"
		if runtime.GOOS == "windows" {
			key = key + ".exe"
		}
		f, err := os.Create("chaosity_temp")
		if err != nil {
			log.Fatalf("There was an error creating temp file %v", err)
		}
		_, err = s3svc.Download(f, &s3.GetObjectInput{
			Bucket: aws.String(updateOptions.UpdateBucket),
			Key:    aws.String(key),
		})
		if err != nil {
			log.Fatalf("There was an error downloading the file %v", err)
		}
		log.Println("Downloaded update file successfully")

		log.Println("Updating binary")

		binary, err := os.Open("chaosity_temp")

		err = update.Apply(binary, update.Options{})

		if err != nil {
			log.Fatalf("Update Failed with error: %v", err)
		}

		log.Println("Update Successful!")

		os.Remove("chaosity_temp")
		log.Println("You are fully updated, to update in the future run again with profile and update flag")

	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

	option := &updateOptions

	updateCmd.Flags().StringVar(&option.Profile, "profile", "", "Profile used to download from S3 - Any standard profile will work (required)")
	updateCmd.Flags().StringVar(&option.UpdateBucket, "bucket", "", "Bucket to download update from - If not known a bucket will be provided to you (required)")
	updateCmd.Flags().StringVar(&option.Region, "region", "us-east-1", "Region specified defaults to us-east-1")

	updateCmd.MarkFlagRequired("profile")
	updateCmd.MarkFlagRequired("bucket")

}

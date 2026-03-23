## Runbook

Log in to https://oyster.tfl.gov.uk/oyster (will probably be send 2FA code via mobile)

Go to the [Journey history](https://oyster.tfl.gov.uk/oyster/journeyHistory.do) and select a custom date range.  Select the earliest `from` date and the latest `to` date possible (you will know because the dates are greyed out).

Scroll to the bottom and click "Download CSV format"

Run `oyster-import rename <downloaded file>` on the downloaded file

Run `oyster-import import <renamed file>` on the renamed file

Run `mv <renamed file> old_csv` to move the file to the archive

Look at the `.env` file and export the `GCP_PROJECT_ID` and `GCP_PUBSUB_TOPIC` environment variables

Run `oyster-import publish`

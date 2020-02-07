This is an action that you can run on a regular basis to clean up outdated
workflow run artifacts before their typical 90 day expiry.

By default, all artifacts are deleted. There are options to only delete big
artifacts, artifacts of a certain age or ones with a specific name.

Suggested usage is in a standalone workflow doc at `.github/workflows/cleanup.yml`
with the following content:

```
name: clean artifacts

on:
  schedule:
    - cron: '0 0 * * *'

jobs:
  clean:
    runs-on: ubuntu-latest
    steps:
      - name: cleanup
        uses: glassechidna/artifact-cleaner@master
        with:
          minimumAge: 86400 # all artifacts at least one day old
```

This cleans all artifacts once a day.

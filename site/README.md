# Docs Site Readme

This directory (site/) on the `main` branch contains all the markdown for the website hosted on a custom domain https://tfx.rocks. 

## How the Site is Built and Published

Changes to the `main` branch in the `site` directory will trigger the [docs-deploy.yml](./../.github/workflows/docs-deploy.yml) Github Action.

This will build the site and push all changes to the `gh-pages` branch, which in turn will trigger the Github action to deploy the site.

The `CNAME` file **must** be present in the `site/docs/` folder for custom domain to be successful.
This file is also copied into the `gh-pages` branch.

## Local Development

```sh
mkdocs serve -f site/mkdocs.yml
```


<p align="center">
  <a href="https://batnoter.com">
    <img src="https://raw.githubusercontent.com/batnoter/batnoter/main/public/logo.svg" width="100">
  </a>

  <p align="center">
    Create and store notes to your git repository!
    <br>
    <a href="https://batnoter.com"><strong>https://batnoter.com</strong></a>
  </p>
</p>

## BatNoter

[![GitHub Workflow Status](https://img.shields.io/github/workflow/status/batnoter/batnoter/Test/main?color=forestgreen)](https://github.com/batnoter/batnoter/actions?query=branch%3Amain)
[![codecov](https://codecov.io/gh/batnoter/batnoter/branch/main/graph/badge.svg?token=P40BDKYDBI)](https://codecov.io/gh/batnoter/batnoter)
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/824dc3f42ddf48f0b99194ea0ef975a7)](https://www.codacy.com/gh/batnoter/batnoter/dashboard?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=batnoter/batnoter&amp;utm_campaign=Badge_Grade)

[BatNoter](https://batnoter.com) is a web application that allows users to store notes in their git repository. This is a frontend project built using mainly react (typescript), redux-toolkit & mui components. [BatNoter API](https://github.com/batnoter/batnoter-api) is the backend implementation of REST APIs which are used by this react app.

<p align="center">
  <kbd><img src="https://raw.githubusercontent.com/batnoter/batnoter/main/demo.gif" alt="BatNoter Demo"/></kbd>
</p>

### Features
-   Login with GitHub.
-   Create, edit, delete, organize & explore notes easily with a nice & clean user interface.
-   Markdown format supported allowing users to add hyperlink, table, headings, code blocks, blockquote... etc inside notes.
-   Editor allows preview of markdown.
-   Quickly copy code from the code section using copy to clipboard button.
-   Store notes directly at the root or use folders to organize them (nesting supported).
-   Explore all the notes from a specific directory with single click.
-   All the notes are stored inside user's github repository.
-   Notes are cached to avoid additional API calls.
-   URLs can be bookmarked.

### Local Development Setup

#### Prerequisites
*   Node.js version `18` or above

#### Start the server
```shell
npm install
npm start
```
This will start the react app in the development mode. Open [http://localhost:3000](http://localhost:3000) to view it in the browser.

#### Run tests
```shell
npm test
```
This will execute all the tests and also prints the code coverage percentage.

### Contribution Guidelines
> Every Contribution Makes a Difference

Read the [Contribution Guidelines](CONTRIBUTING.md) before you contribute.

### Contributors
Thanks goes to these wonderful people 🎉

[![](https://opencollective.com/batnoter/contributors.svg?width=890&button=false)](https://github.com/batnoter/batnoter/graphs/contributors)

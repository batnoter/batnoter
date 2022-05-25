<p align="center">
  <a href="https://gitnoter.com">
    <img src="https://raw.githubusercontent.com/git-noter/gitnoter/main/public/logo.svg" width="100">
  </a>

  <p align="center">
    Create and store notes to your git repository!
    <br>
    <a href="https://gitnoter.com"><strong>https://gitnoter.com</strong></a>
  </p>
</p>

## GitNoter

[![GitHub Workflow Status](https://img.shields.io/github/workflow/status/git-noter/gitnoter/Test/main?color=forestgreen)](https://github.com/git-noter/gitnoter/actions?query=branch%3Amain)
[![codecov](https://codecov.io/gh/git-noter/gitnoter/branch/main/graph/badge.svg?token=P40BDKYDBI)](https://codecov.io/gh/git-noter/gitnoter)
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/276bf59cb3ba4249863bbdb8c290fe14)](https://www.codacy.com/gh/git-noter/gitnoter/dashboard?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=git-noter/gitnoter&amp;utm_campaign=Badge_Grade)

[GitNoter](https://gitnoter.com) is a web application that allows users to store notes in their git repository. This is a frontend project built using mainly react (typescript), redux-toolkit & mui components. [GitNoter API](https://github.com/git-noter/gitnoter-api) is the backend implementation of REST APIs which are used by this react app.

<p align="center">
  <img src="https://raw.githubusercontent.com/git-noter/gitnoter/main/public/demo/demo-gitnoter-720p.gif" alt="GitNoter Demo"/>
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

[![](https://opencollective.com/gitnoter/contributors.svg?width=890&button=false)](https://github.com/git-noter/gitnoter/graphs/contributors)

Contributing Guidelines
=======================
Contributing to this project does not only help it grow but it shows us your highly valued interests to use your talent to contribute to an open source project like this and I am very eager to welcome you to contribute! As a matter of fact, this is what any Open Source project is all about. Before we work on with any contributions, the purpose of this contributing document is to guide you how your contribution(s) should properly align to how this project is intentionally designed and maintained based on certain guidelines which will be discussed below. We ask you to read this document in its entirety.

Remember, these are guidelines, not rules. The guidelines are set to create a "communication" between the original author and you. Sometimes, we may have different interpretations of the guidelines discussed here but you have the freedom to use your best judgment. It can be complicating but should your pull requests be ever put into thorough review, explanations will be provided as comments. Don't worry! Contributions are almost guaranteed as long as we work on fixing them if there are any issues.

---

Contributing License Agreement
==============================
Before contributing and your pull requests to be accepted, we ask you to electronically sign the [Contributing License Agreement](https://cla-assistant.io/rrborja/minesweeper) for the protection of everyone who made this project successfully grow.

For corporations who want to contribute, please email to inquiry@brute.io

Code of Conduct
===============
All contributing members become community members automatically and are expected to adhere to our [code of conduct](https://github.com/rrborja/minesweeper/blob/master/CODE_OF_CONDUCT.md).

License
=======
This project Minesweeper API is released under the [GNU General Public License](https://www.gnu.org/licenses/old-licenses/gpl-2.0.en.html), either version 2.0 or later versions.

Report a Bug
============
Bug fixes are really not fun to play with, compared to creating functionalities. However, we must always recognize that, in time, there will be bugs. Bug fixes help close the missing gaps of the software that made them vulnerable for exploits and opportunistic unusual behaviors. Let us know if there is a bug by creating an issue in the Minesweeper GitHub repository.

Creating an Issue
=================
Before submitting an issue, ensure that:
* at the time of the issue present in the project, the source code must be at its latest revision.
* the issue has not been described in the [issue tracker](https://github.com/rrborja/minesweeper/issues) yet.

Suggest a Functionality
=======================
Adding a functionality is not only a fun task but it is always a serious task to become involved with. As such, there will be a lot of discipline that both of us must expect. The following guidelines will help us achieve the best overall outcome:
1. Every feature development must be coupled with test-driven development.  
   * If the functionality is not yet existed, test cases must be created first and ensure the test fails at first. The reason behind is that when a test fails after performing a testing scenario, we determine that the functionality has not been  implemented or, better yet, the code for that functionality is not in the codebase yet. When it does fail, it will now become your opportunity to create (a.k.a refactor) that functionality to the actual codebase. The second test then must pass.
2. Your naming conventions must be properly aligned to this [guideline](https://golang.org/doc/effective_go.html).
3. Your exported methods, functions, types, variables and constants must have their own Godoc-style comments. Comments should be meaningful and helpful in order to provide other contributors enough idea how your code works.
4. Your changes must pass code vetting frameworks. The frameworks used are indicated by the top of the [README.md](https://github.com/rrborja/minesweeper/blob/master/README.md), in this case, travis-ci, goreportcard, coveralls, etc.
5. Continuously communicate with us what is happening during your work. It helps us know that ideas still keep flowing.

Opening a Pull Request
======================
[Fork](https://help.github.com/articles/fork-a-repo/) `rrborja/minesweeper`, commit your changes, and [open a pull request](https://github.com/rrborja/minesweeper/compare).

I will be notified immediately when you have submited the pull request. During this process, your submission will be reviewed and vetted with most automated code review frameworks. We ask for your patience and you will be notified accordingly.

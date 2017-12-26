# Minesweeper API

[![GoDoc](https://godoc.org/github.com/rrborja/minesweeper?status.svg)](https://godoc.org/github.com/rrborja/minesweeper)
[![License: GPL v2](https://img.shields.io/badge/License-GPL%20v2-blue.svg)](https://www.gnu.org/licenses/old-licenses/gpl-2.0.en.html)  
[![Build Status](https://travis-ci.org/rrborja/minesweeper.svg?branch=master)](https://travis-ci.org/rrborja/minesweeper)
[![codecov](https://codecov.io/gh/rrborja/minesweeper/branch/master/graph/badge.svg)](https://codecov.io/gh/rrborja/minesweeper)
[![Go Report Card](https://goreportcard.com/badge/github.com/rrborja/minesweeper)](https://goreportcard.com/report/github.com/rrborja/minesweeper)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/cec50a1b138e4e7789a7ffb0e61432e4)](https://www.codacy.com/app/rrborja/minesweeper?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=rrborja/minesweeper&amp;utm_campaign=Badge_Grade)
[![Maintainability](https://api.codeclimate.com/v1/badges/3c3d7b7aef3cec7c7ef5/maintainability)](https://codeclimate.com/github/rrborja/minesweeper/maintainability)

The Minesweeper API is the Application Programming Interface of the minesweeper game. The game's logic is embedded in this API itself with zero third-party libraries involved. 

Although, there is much no use of this game as a library, you may find this useful for educational purposes because it contains the implementation of the minesweeper game backed with coding vetting frameworks and test-driven development. The idea of this project is to provide you how this project integrates with DevOps, how it can be consumed with REST API, and how the code can achieve high code quality.

The best part: this library contains 100% of all source codes written in Go.

---

* [Install](#install)
* [Usage](#usage)
  * [Creating the Instance](#creating-the-instance)
  * [Setting the Grid](#setting-the-grid)
  * [Setting the Difficulty](#setting-the-difficulty)
  * [Start the Game](#start-the-game)
  * [Visit a Cell](#visit-a-cell)
  * [Flag a cell](#flag-a-cell)
* [Example](#example)
* [TODO](#todo)
* [License](#license)
* [Contributing](#contributing)


---

Install
=======

`go get -u github.com/rrborja/minesweeper`

Usage
=====

### Creating the Instance
Create the instance of the game by calling `minesweeper.NewGame()` method. The method can accept an arbitrary number of arguments but only one argument can be processed and the rest are ignored. The method accepts the `Grid` argument to set the board's size. For example, `Grid{Width: 15, Height: 10}`

> **Setting the Grid:**  
> If the `Grid` is not provided in the `NewGame()` arguments, you must explicitly provide the board's `Grid` by calling `SetGrid()` of the game's instance.

`NewGame()` returns two values: the instance itself and the event handler. The instance is the instance of the `Minesweeper` interface that has methods as use cases to solve a minesweeper game. The event handler is a buffered channel that you can use to create a separate goroutine and listen for game events. Such events are `minesweeper.Win` and `minesweeper.Lose`.

### Setting the Difficulty
Set the difficulty of the game by calling `SetDifficulty()` of the game's instance. Values accepted by this method as arguments are `minesweeper.Easy`, `minesweeper.Medium` and `minesweeper.Hard`.

### Start the Game
Call the `Play()` of the game's instance to generate the location of mines and to start the game. You may encounter errors such as `UnspecifiedGridError` if no grid is set, `UnspecifiedDifficultyError` if no difficulty is set, and `GameAlreadyStartedError` if the `Play()` has already been called twice or more.

### Visit a Cell
Call the `Visit()` of the game's instance to visit the cell. The method will accept two arguments of the type `int` which are represented by the xy-coordinate of the game's board to which the location of the cell in the board is to be visited.  
> **The method will return two values:**
> - a slice of cells that has been visited by the player and by the game
> - the error `ExplodedError` indicating whether a visited cell reveals a mine.

When you receive the `ExplodedError` error, the game ends. Calling `Visit()` after a game ended may produce a nil pointer dereference panic. It does happen due to that the game is over, the game's instance is automatically garbage collected.

The slice being returned may contain multiple values. If it's a single element slice, that means the visited cell is a warning number indicating the number of mines neighbored to the visited cell. If it's a multiple element slice, that means the visited cell is an unknown number and neighboring cells are recursively visited until any numbered cell is reached. The first element of the slice will always be the original player's visited cell.

### Flag a cell
Call `Flag()` of the game's instance to mark an unprobed cell. Doing this will prevent the `Visit()` method from visiting the marked cell.

> **Pro tip:**  
> - If you visit an already visited numbered cell again and the number of neighboring cells that have been flagged equals to the number of the visited cell, the game will automatically visit all unprobed neighbored cells by returning the slice containing those cells. Just make the players ensure that they have correctly marked the cells deduced that they have mines, otherwise, the game will end if the cell is incorrectly marked.

Example
=======

```go
package main

import (
    "fmt"
    "os"

    "github.com/rrborja/minesweeper"
)

func Listen(event Event) {
    select event {
    case <- minesweeper.Win:
        fmt.Println("You won the game!")
    case <- minesweeper.Lose:
        fmt.Println("Game over!")
    default:
        panic("Unexpected event")
    }
    os.Exit(0)
}

func main() {
    game, event := minesweeper.NewGame(Grid{Width: 15, Height: 10})
    game.SetDifficulty(Easy)
    game.Play()

    go Listen(event)

    reader := bufio.NewReader(os.Stdin)
    for {
        var x, y int
        fmt.Println("Enter the X-coordinate to visit a cell: ")
        fmt.Scan(&x)
        fmt.Println("Enter the Y-coordinate to visit a cell: ")
        fmt.Scan(&y)

        visitedCells, err := game.Visit(x, y)
        if err != nil {
            fmt.Printf("Mine is revealed at %v:%v", x, y)
        }

        for _, cell := range visitedCells {
            fmt.Print("[%v:%v] ", cell.X(), cell.Y())
        }
    }
}
```

TODO
====
1. Print the game's visual state of the board in console
2. Provide a way to allow the game's API to be used for REST API

License
=======

This project Minesweeper API is released under the [GNU General Public License v2.0](https://www.gnu.org/licenses/old-licenses/gpl-2.0.en.html).

GPG Verified
============

Commits are signed by verified PGPs

Contributing
============

Please see the [CONTRIBUTING.md](https://github.com/rrborja/minesweeper/blob/master/CONTRIBUTING.md) file for details.
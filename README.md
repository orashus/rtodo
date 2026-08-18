![preview img](./preview.png)

## Index

1. [rtodo](#rtodo)
2. [Install](#install)
3. [Storage](#storage)
4. [Todo shape](#todo-shape)
5. [Usage](#usage)
6. [Commands](#commands)
7. [Flags](#flags)
8. [Examples](#examples)

## rtodo

`rtodo` is a small command-line task tracker. It stores tasks as todos in a JSON file on disk in your `/tmp/` dir and prints them as a table.

The base code with the initial commits of this project live on my go-sandbox repo ([github.com/RashJrEdmund/go-sandbox](https://github.com/RashJrEdmund/go-sandbox)).
This project started out as a means of learning file handling.

```bash
rtodo <command> [inputs...] [flags]
```

## Install

Requires Go. Install from the [main repository](https://github.com/orashus/rtodo). This puts `rtodo` on your `PATH` (`$GOPATH/bin` or `$HOME/go/bin`):

```bash
go install github.com/orashus/rtodo@latest
```

From this learning copy instead:

```bash
go install .
```

If you run it with no command:

```text
rtodo  version 1.0.0
Please provide a command
```



## Storage

Todos are loaded from and saved to `test-todos.json` in this directory.

If the file does not exist, the list is treated as empty. Saves write pretty-printed JSON (2-space indent). `clear` deletes the file.

The sample file currently has one todo. The examples below use it:


| ID      | Title           | Completed | Created at                            |
| ------- | --------------- | --------- | ------------------------------------- |
| `zAMzK` | Buy a new house | `false`   | `2026-08-18T16:34:53.480503493+01:00` |




## Todo shape

Each todo is a JSON object:

```json
{
  "id": "zAMzK",
  "title": "Buy a new house",
  "completed": false,
  "created_at": "2026-08-18T16:34:53.480503493+01:00"
}
```


| Field        | Type      | Notes                                                                                            |
| ------------ | --------- | ------------------------------------------------------------------------------------------------ |
| `id`         | string    | 5-character unique id (`a-zA-Z0-9_`)                                                             |
| `title`      | string    | The todo text you pass to `add`, or the new text you pass to `update`                            |
| `completed`  | bool      | New todos start as `false`; `complete` / `done` / `finish` / `check` / `mark` set this to `true` |
| `created_at` | timestamp | Set to `time.Now()` when added                                                                   |


When printed, the table looks like:

```text
ID         Created At                   Completed  Title
---------------------------------------------------------------------------------------------------
zAMzK      2026-08-18T16:34:53+01:00    false      Buy a new house
```



## Usage

Command is the first argument. Everything after it (that is not a flag) is input. Flags can appear anywhere in the remaining args, as `-short` or `--long`.

`add`, `rm` / `remove` / `delete`, and `complete` / `done` / `finish` / `check` / `mark` all take **one or more** inputs. `update` takes **two**: an id, then a new title.

Quote any title that contains spaces. Without quotes, each word is a separate input (`rtodo add Buy a house` adds three todos).

```bash
rtodo list
rtodo add "Buy a new house"
rtodo add "Buy a new house" "Walk the dog"
rtodo done zAMzK --print
rtodo done zAMzK abc12
rtodo update zAMzK "Sell the house"
rtodo rmc --print
rtodo rm zAMzK
rtodo rm zAMzK abc12
rtodo list --completed
rtodo --version
```

IDs are generated automatically on `add`. Look them up with `list` (or in the JSON file) before `rm`, `complete`, `update`, and the other id-based commands.

## Commands



### `list`

Prints every todo in the store.

```bash
rtodo list
```

With `--completed` / `-c`, only completed todos are shown.

```bash
rtodo list --completed
```

With the sample file this is empty, because `zAMzK` has `"completed": false`.

### `add <title> [title...]`

Appends one incomplete todo per title. Requires at least one title.

```bash
rtodo add "Buy a new house"
rtodo add "Buy a new house" "Walk the dog"
```

Success:

```text
---------------------------------------------------------------------------------------------------
Todo(s) added successfully
```

If no title is given:

```text
Please provide a title to add
```

Pass `--print` / `-p` or `--list` / `-l` to print the full list after adding.

### `complete` / `done` / `finish` / `check` / `mark` `<id> [id...]`

Marks those todos as complete (`completed: true`). All five command names do the same thing. Requires at least one id.

```bash
rtodo complete zAMzK
rtodo done zAMzK abc12
rtodo finish zAMzK
rtodo check zAMzK
rtodo mark zAMzK
```

Success (one id):

```text
---------------------------------------------------------------------------------------------------
Todo(s) with ID(s) '[zAMzK]' marked as complete
```

Success (several ids):

```text
---------------------------------------------------------------------------------------------------
Todo(s) with ID(s) '[zAMzK abc12]' marked as complete
```

If no id is given:

```text
Please provide an ID to mark as complete
```

If none of the ids are in the store:

```text
---------------------------------------------------------------------------------------------------
Todo(s) not found
```

If at least one id matches, those todos are marked complete even if other ids in the same command are missing.

Pass `--print` / `-p` or `--list` / `-l` to print the list afterwards. After `done zAMzK --print`:

```text
ID         Created At                   Completed  Title
---------------------------------------------------------------------------------------------------
zAMzK      2026-08-18T16:34:53+01:00    true       Buy a new house
```

### `update <id> <title>`

Changes the title of the todo with that id. Requires both arguments. Quote the title if it contains spaces; only the first two inputs are used (`id`, then `title`).

```bash
rtodo update zAMzK "Sell the house"
```

Success:

```text
---------------------------------------------------------------------------------------------------
Todo with ID 'zAMzK' updated successfully
```

If the id or title is missing:

```text
Please provide an ID and a title to update
```

If the id is not in the store:

```text
---------------------------------------------------------------------------------------------------
Todo not found
```

Pass `--print` / `-p` or `--list` / `-l` to print the list afterwards.

### `rm` / `remove` / `delete` `<id> [id...]`

Removes the todos with those ids. All three command names do the same thing. Requires at least one id.

```bash
rtodo rm zAMzK
rtodo rm zAMzK abc12
rtodo remove zAMzK
rtodo delete zAMzK
```

Success (one id):

```text
---------------------------------------------------------------------------------------------------
Todo(s) with ID(s) '[zAMzK]' removed successfully
```

Success (several ids):

```text
---------------------------------------------------------------------------------------------------
Todo(s) with ID(s) '[zAMzK abc12]' removed successfully
```

If no id is given:

```text
Please provide an ID to remove
```

If none of the ids are in the store:

```text
---------------------------------------------------------------------------------------------------
Todo(s) not found
```

If at least one id matches, those todos are removed even if other ids in the same command are missing.

Pass `--print` / `-p` or `--list` / `-l` to print the remaining list after removal.

### `rmc`

Removes every completed todo. Incomplete todos are left in place.

```bash
rtodo rmc
```

Success:

```text
---------------------------------------------------------------------------------------------------
Completed todos removed successfully
```

If `zAMzK` is still incomplete, `rmc` leaves it. After `done zAMzK` then `rmc --print`, the list is empty.

### `clear`

Deletes the todo file. If there is nothing to delete:

```text
---------------------------------------------------------------------------------------------------
No todos to clear
```

Otherwise:

```text
---------------------------------------------------------------------------------------------------
Todos cleared successfully
```

Pass `--print` / `-p` or `--list` / `-l` to print the (now empty) table afterwards.

### No command / `--version` / `-v`

`--version` / `-v` prints:

```text
rtodo  version 1.0.0
```

and exits without running a command.

An unknown command prints `Invalid command`. With `--print` / `--list` it still dumps the current list.

## Flags

Short and long forms are equivalent. Parsed flags are stripped before the command and input are read.


| Short | Long          | Effect                                                                                                  |
| ----- | ------------- | ------------------------------------------------------------------------------------------------------- |
| `-p`  | `--print`     | After `add`, `rm`, `rmc`, `clear`, `update`, or `complete`/`done`/`finish`/`check`/`mark`, print the current list |
| `-l`  | `--list`      | Same as `--print`                                                                                       |
| `-c`  | `--completed` | With `list`, show only completed todos                                                                  |
| `-v`  | `--version`   | Print `rtodo version 1.0.0` and exit                                                                    |
| `-h`  | `--help`      | Recognized, but no help text is printed yet                                                             |




## Examples

Using the sample todo from `test-todos.json`.

**List everything**

```bash
rtodo list
```

```text
ID         Created At                   Completed  Title
---------------------------------------------------------------------------------------------------
zAMzK      2026-08-18T16:34:53+01:00    false      Buy a new house
```

**Add another todo, then print**

```bash
rtodo add "Walk the dog" --print
```

**Add several todos at once**

```bash
rtodo add "Walk the dog" "Buy milk"
```

**Rename a todo**

```bash
rtodo update zAMzK "Sell the house" --print
```

```text
---------------------------------------------------------------------------------------------------
Todo with ID 'zAMzK' updated successfully
ID         Created At                   Completed  Title
---------------------------------------------------------------------------------------------------
zAMzK      2026-08-18T16:34:53+01:00    false      Sell the house
```

**Mark "Buy a new house" complete**

```bash
rtodo done zAMzK --print
```

```text
---------------------------------------------------------------------------------------------------
Todo(s) with ID(s) '[zAMzK]' marked as complete
ID         Created At                   Completed  Title
---------------------------------------------------------------------------------------------------
zAMzK      2026-08-18T16:34:53+01:00    true       Buy a new house
```

**Mark several todos complete**

```bash
rtodo done zAMzK abc12
```

**List only completed todos**

```bash
rtodo list -c
```

After the `done` command above, that prints `zAMzK`. Before it, the list is empty.

**Remove completed todos**

```bash
rtodo rmc --print
```

That drops `zAMzK` once it is complete, and keeps any incomplete items.

**Remove by id**

```bash
rtodo rm zAMzK --print
```

**Remove several todos by id**

```bash
rtodo rm zAMzK abc12 --print
```

**Wipe the list**

```bash
rtodo clear
```
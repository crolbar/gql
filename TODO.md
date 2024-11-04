-   [x] table scrollbar
-   [x] show that there are more rows to the bottom/top
-   [x] fix crash on cursor move in empty table
-   [x] fix horizontal scroll not worky
-   [x] auth fix errors with switching database from a keybing
-   [x] don't hardcode width's and height's
    -   ~~[ ] avg column width~~
-   [x] describe ~~pane~~ tab
-   [x] curr table info (columns rows)
-   [x] show curr user
-   [x] help view
-   [x] if an keypress is handled by the table update don't handle it from the pane's update (eg: esc)
-   [x] esc for cancel in auth

-   [x] errors (eg: no permisions to access table)

-   [x] filter (where clause)
    -   [x] set textinput width
-   [x] update
    -   [x] cell
    -   [x] table name
-   [ ] ~~insert~~
-   [ ] delete / drop
    -   [x] drop db
    -   [x] drop table
    -   [x] delete row
    -   [x] delete selected rows
    -   [ ] delete selected columuns
-   [ ] ~~create~~

-   [x] dialog

    -   [x] help msg
    -   [x] err msg
    -   [ ] ~~Table help msg ?~~
    -   [x] string (non confirmation) dialog for updating names / columns

-   [x] send a custom sql query
    -   [ ] show result

-   postgresql support fixes
    -   [x] fix update.go on line 26
    -   [x] fix where caluse builder for postgress (timestamp can't be matched with an string)
    -   [ ] database switching in postgress ?
    -   [ ] use exec instead of query on non queries
    -   [ ] close db connection on uri change

-   [x] when dialog is selected don't process any char buttons
-   [x] fix crash on none tables in db
-   [x] fix crashes on empty main table
-   [ ] refresh dbs and db tables button
-   [ ] more to help menu
-   [ ] fix binary data

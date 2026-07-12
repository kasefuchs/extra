# Teardown promo

for some dumb reason Teardown uses plain http for featured mods stuff

so yeah... you can just hijack it

[see it in action](https://social.floof.fans/@kasefuchs/statuses/01KMXHYNQ9T646PQV5XD9DD616)

### how it works

1. the game asks `http://promo.teardowngame.com/promo.php?version=x.x.x` what to download
2. server replies with a raw link to an archive: `http://promo.teardowngame.com/YY-MM-DD.zip`
3. then game downloads and extracts it. all over plain http

### how to use

you'll need `go-task`, `zip`, and `python3`

1. build the payload:

   ```sh
   go-task
   ```

2. start the fake server:

   ```sh
   sudo go-task serve
   ```

3. edit your `/etc/hosts`:

   ```text
   127.0.0.1 promo.teardowngame.com
   ```

4. launch the game and check the featured mods page

## notes

- works because of plain http lol
- breaks the moment they switch to https
- dumb poc but funny

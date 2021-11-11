curl -v --request POST \
  --url 'https://login.devhost.dev/token' \
  --header 'content-type: application/x-www-form-urlencoded' \
  --data grant_type=password \
  --data 'client_id=2OCXpewmSmS8qhWKooXECg' \
  --data username=geoffcake@gmail.com \
  --data 'password=Example#1' \
  | jq


curl -v --request POST \
  --url 'https://login.devhost.dev/token' \
  --header 'content-type: application/x-www-form-urlencoded' \
  --data grant_type=client_credentials \
  --data 'client_id=2OCXpewmSmS8qhWKooXECg' \
  | jq

    # --data client_secret=YOUR_CLIENT_SECRET \

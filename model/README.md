A context has many clients, users, groups, roles, actions and ratelimits.
a user can belong to many groups
a group has many users
a role has many actions

there is a link table between user/group pairs, and roles.
Each user may have different roles in each group

role->action mapping is global per context
action rate limit is global per context

a user has many auth methods

for ease of storage, auth codes are issues in encrypted jwts, so they don't need storage.
when a browser auths, a session cookie is placed.  jwt of that session cookie is stored.
when an access token is issued, its jti, refresh token, and the session jti are stored.
session jti is stored in the access token

actions are grant only, no deny.

actions are ltrees.  being granted the root of the tree is being granted "root", and is exactly what it sounds like


Pbkdf(Secret+time mod 5min), gives password. Can use sliding window like totp, but for encryption key. 
the 5min part should be configurable
one function that takes the secret, and the time, and then does the code, and another that handles the offset bits, and can return a list of valid codes for a specific time window.
Might also just use a simple list of possible secrets stored in the config
Or, just have a list of valid keys in the config. First key is used for encryption, and the others are tried for decryption. 

Put the browser/session token in with auth code data, then can tie everything together when issuing access token. 
basically, we just need to put the associated authed browser session into the encrypted auth code.  That way, when the code turns back into an access token,
we know which browser session it should be associsated with.  if issuing an access token directly, we have the browser session handy, and can just directly use it from there.
We'll need an "access token context" table, since a single context can be used to have multiple valid access tokens, by renewal and whatnot.  Will want to make sure that
the bits and bops are all lined up right there.

Should have a token util file, with helpers for serialization, sign, encrypt and the like. 

Consolidate Janus notes as well


seperate user from context.  When doing a signup, detect that the user exists, and when passing the new user data to the signup app, just let it create the group.
when authenticating, we'll know the client, and hence can find the appropriate user.

can I make the refresh token something that I can just use to lookup the access token?  like make the refresh token be somehow meaningfully derived from the access token, or a jwt that has it signed inside it?  that way I don't need to store it, but can verify that it's correct?
I'm thinking use an encrypted jwt, that way I can just unseal it and have the data, no storage needed.


Put the different http handlers under different namespaces, so that very little loves at the top level.  This will make it easier to organize things. 

If a user hits the main page, log them in, and then present them with a list of every app they could sign into, likely derived from their group memberships. 

Also, gravatars everywhere. 

For refresh token, use an hmac of the access token I'd? Will need an hmac helper... -- nope, need to be able to do access token lookup from the refresh token.  encrypt it!

webauthn support should be easy-ish, and kinda neat.  Should be able to make it work with just a small "biometric" button on the login page, and the info in the handshake can be used to identify the user
implies that the auth-check method should be able to understand that it might be doing a lookup on something other than email

Support enrollment protocol?if hit main page, show grid of contexts (available by group), and drill to grid of apps. If only one, start with apps. 

Need a link between users and contexts/clients that they've "joined"?  Or is that over designing?

Can we represent refresh tokens with jwe, and find a way to not store jti in db either?  Just store the session per client, and use encryption to de-auth?
Session has a signing key, and that signs a value in the refresh token?

What if we just stick the session code in the access token?  then we can also just use a jwe of the session data as the refresh token.  Would also need to stick the jti of the access token in there, but that's no biggie

Access Token has a sessid claim, and a browser id claim.  That way it can be known what browser and session an access token came from.  The refresh token will have the sesssion id and the access token id.  That way it can be confirmed that the refresh token matches with the given access token, and it can also be revoked by revoking the session id

Need to add a lot of recovery logic for if things get out of sync.  Can have sessions just falling over all the time

consider abandoning context notion.  Too complicated, not enough gain.  Straight forward to re-add if it becomes wanted later

Need to pass the context into identification functions, and can use that to identify the cookie in question. Need to put the context code into the I'd token, so that we know the cookie name matches the I'd token

Might also look at finding a better way to tie context to user. Not sure matters. 

Also need to validate token expiration times.  

Would it work to store user ssh keys, as a fun diversion?  Use AuthorizedKeysCommand as a way to validate pup keys for users, and give permissions for different users to login as different users?  Could be neat. 
Should definitely do that. Could work as a neat way to authenticate users even across web requests.   With some trickery, would be possible to make it divert auth to fingerprint reader?
Need to link the key, identity, and what auth user it can login as. Also indicate if a user can use key to login to anyone. 

Look at browser compiled react/Babel stuff. https://reactjs.org/docs/add-react-to-a-website.html

Client API should have way to manage user group memberships, and need to think about ways to handle checking for user groups, for asking if users are in the same group.  
But maybe just need to say if a group is in a context, that context can manage group memberships, and possibly just tracks that in the client app. 

for the webauthn session bits, can store the bits relating to doing the handshake in a jwe, which saves us from having to do more db stuff
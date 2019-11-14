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
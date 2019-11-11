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
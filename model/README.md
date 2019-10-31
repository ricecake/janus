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
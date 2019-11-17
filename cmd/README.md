Need to restructure the commands, so that instead of create/update/delete being toplevel, 
we make user/client/context etc top level, with create/update/delete/details underneath.

Then the crud commands can be kinda generic-ish, and can use pre/post execute hooks on the object level to 
setup the flags and create the object in question, with them all having the same, roughly, interface for how you crud the object.

so an interface for the basic objects that mandates basic crud, and that way the commands should almost write themselves

can nestle the entire thing under an `admin` command, to keep the top level from getting too full
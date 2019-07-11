converter accepts callbacks from bellbox ignorant webhooks
 - converts from a source encoding to a json message
 - handles authenticating as a sender, suppressing errors
 - configured with a local path which must be called by the remote service

converter is statically configured with json files.

 - translations/ contains files which contain translation rulesets. each ruleset is named after their file and the individual rule name. this allows for optional formatting
 - config.json contains a list of paths mapped to translation rulesets
 - MAYBE: config/ contains other files to merge with config.json

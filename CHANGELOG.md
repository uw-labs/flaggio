<a name="unreleased"></a>
## [Unreleased]


<a name="v0.4.0"></a>
## [v0.4.0] - 2020-04-13
### Bug Fixes
- nils from context should be treated as invalid ([#33](https://github.com/victorkt/flaggio/issues/33))

### Dependencies
- fix package-lock.json ([#53](https://github.com/victorkt/flaggio/issues/53))
- bump [@material](https://github.com/material)-ui/styles from 4.9.6 to 4.9.10 in /web ([#51](https://github.com/victorkt/flaggio/issues/51))
- bump [@material](https://github.com/material)-ui/core from 4.9.9 to 4.9.10 in /web ([#49](https://github.com/victorkt/flaggio/issues/49))
- bump go.mongodb.org/mongo-driver from 1.3.1 to 1.3.2 ([#52](https://github.com/victorkt/flaggio/issues/52))
- bump github.com/go-chi/chi from 4.0.4+incompatible to 4.1.0+incompatible ([#35](https://github.com/victorkt/flaggio/issues/35))

### Documentation
- add CHANGELOG

### Features
- new variant is consistent with previous ones ([#46](https://github.com/victorkt/flaggio/issues/46))
- implement distribution rollout on frontend ([#43](https://github.com/victorkt/flaggio/issues/43))
- case insensitive sort for find all flags ([#44](https://github.com/victorkt/flaggio/issues/44))
- add 'rows per page' = All ([#42](https://github.com/victorkt/flaggio/issues/42))
- add tracing support ([#41](https://github.com/victorkt/flaggio/issues/41))
- caching is optional ([#40](https://github.com/victorkt/flaggio/issues/40))
- add redis cache support ([#34](https://github.com/victorkt/flaggio/issues/34))

### Performance Improvements
- limit amount of variants/rules/constraints ([#45](https://github.com/victorkt/flaggio/issues/45))


<a name="v0.3.0"></a>
## [v0.3.0] - 2020-04-01
### Bug Fixes
- evaluateAll doesn't return nulls anymore ([#31](https://github.com/victorkt/flaggio/issues/31))
- remove default variant refs in newFlag func ([#30](https://github.com/victorkt/flaggio/issues/30))

### Features
- add pagination to flags table ([#29](https://github.com/victorkt/flaggio/issues/29))
- add flag search functionality ([#28](https://github.com/victorkt/flaggio/issues/28))
- default variants for new flags ([#25](https://github.com/victorkt/flaggio/issues/25))


<a name="v0.2.0"></a>
## v0.2.0 - 2020-03-28
### Bug Fixes
- gqlgen new server api ([#20](https://github.com/victorkt/flaggio/issues/20))

### Continuous Integration
- only release on tags
- build & release ([#24](https://github.com/victorkt/flaggio/issues/24))

### Dependencies
- bump clientip ([#23](https://github.com/victorkt/flaggio/issues/23))

### Work in Progress
- new cmd


[Unreleased]: https://github.com/victorkt/flaggio/compare/v0.4.0...HEAD
[v0.4.0]: https://github.com/victorkt/flaggio/compare/v0.3.0...v0.4.0
[v0.3.0]: https://github.com/victorkt/flaggio/compare/v0.2.0...v0.3.0

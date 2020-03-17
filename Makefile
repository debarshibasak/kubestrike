all: revive
revive:
	revive -config ./.project/revive.toml -formatter friendly -exclude ./vendor/... ./...
release:
	./.project/semver.sh -p
release-major:
	./.project/semver.sh -M
release-minor:
	./.project/semver.sh -m

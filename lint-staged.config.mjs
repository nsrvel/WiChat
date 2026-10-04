/** @type {import('lint-staged').Configuration} */
export default {
	"frontend/web-client/**/*.{ts,tsx,mjs,cjs}": [
		"pnpm --filter web-client exec eslint --fix --max-warnings 0",
	],
	"frontend/web-client/**/*.{json,css,md,yml,yaml}": [
		"pnpm --filter web-client exec prettier --write",
	],
};

const glob = require('glob')
const path = require("path");
const webpack = require("webpack");

let files = glob.sync('./src/pages/**/index.jsx').reduce((acc, path) => {
	const entry = path.replace('./src/', '').replace('/index.jsx', '')
	acc[entry] = path.replace('/src', '');
	return acc
}, {});

module.exports = {
	entry: files,
	output: {
		filename: './build/[name].js',
		path: path.resolve(__dirname)
	},
	context: path.resolve(__dirname, 'src/'),
	resolve: {
		extensions: ["*", ".js", ".jsx"],
		modules: [path.resolve('./node_modules')],
		alias: {
			Component: path.resolve(__dirname, 'src/components/')
		}
	},
	optimization: {
		minimize: true,
		usedExports: true,
	},
	mode: "production",
	module: {
		rules: [
			{
				test: /\.(js|jsx)$/,
				exclude: /(node_modules|bower_components)/,
				loader: "babel-loader",
				options: { presets: ["@babel/env"] }
			},
			{
				test: /\.css$/,
				use: ["style-loader", "css-loader"]
			}
		]
	},
};

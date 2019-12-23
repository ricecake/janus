const glob = require('glob')
const path = require("path");


let libs = glob.sync('./src/libs/**/index.js').reduce((acc, path) => {
	const entry = path.replace('./src/', '').replace('/index.js', '')
	acc[entry] = path;
	return acc
}, {});

let files = glob.sync('./src/pages/**/index.jsx').reduce((acc, path) => {
	const entry = path.replace('./src/', '').replace('/index.jsx', '')
	acc[entry] = path;
	return acc
}, libs);

module.exports = {
	entry: files,
	output: {
		filename: './build/[name].js',
		path: path.resolve(__dirname),
	},
	context: path.resolve(__dirname),
	resolve: {
		extensions: ["*", ".js", ".jsx"],
		modules: [path.resolve('./node_modules')],
		alias: {
			Component: path.resolve(__dirname, 'src/components/'),
			Include: path.resolve(__dirname, 'src/includes/'),
		}
	},
	optimization: {
		minimize: true,
		usedExports: true,
		runtimeChunk: 'single',
		splitChunks: {
			cacheGroups: {
				vendor: {
					test: /[\\/]node_modules[\\/]/,
					name: 'vendors',
					chunks: 'all',
				},
			},
		},
	},
	mode: "production",
	module: {
		rules: [
			{
				test: /\.(js|jsx)$/,
				exclude: /(node_modules|bower_components)/,
				loader: "babel-loader",
				options: {
					presets: ["@babel/env"],
					"plugins": ["minify-dead-code-elimination"],
				}
			},
			{
				test: /\.css$/,
				use: ["style-loader", "css-loader"]
			}
		]
	},
};

console.log(module.exports);
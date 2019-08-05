const path = require('path');
const MiniCssExtractPlugin = require('mini-css-extract-plugin');
const webpack = require('webpack');

module.exports = {
	mode: 'production',
	resolve: {
		alias: {
			jquery: "jquery/src/jquery",
		}
	},
	plugins: [
		new MiniCssExtractPlugin({
			filename: 'css/[name].css',
			chunkFilename: '[id].css',
		}),
	],
	entry: {
		main: './src/index.js',
	},
	output: {
		filename: 'js/[name].js',
		path: path.resolve(__dirname, 'dist')
	},

	/*optimization: {
		splitChunks: {
			chunks: 'all'
		}
	},*/

	module: {
		rules: [
			{
				test: /\.js$/,
				exclude: /(node_modules)/,
				use: {
					loader: 'babel-loader',
					options: {
						presets: ['@babel/preset-env']
					}
				}
			},
			{
			test: /\.(sa|sc|c)ss$/,
			use: [
				MiniCssExtractPlugin.loader,
				{
					loader: 'css-loader',
				}, {
					loader: 'postcss-loader',
					options: {
						plugins: function() {
							return [
								require('autoprefixer')
							];
						}
					}
				}, {
					loader: 'sass-loader'
				}]
		}]
	},
};



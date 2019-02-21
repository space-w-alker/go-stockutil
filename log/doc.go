/*
Package log package provides convenient and flexible utilities for logging messages.

Overview

Logging is done by calling functions in this package that correspond to the
severity of the log message being output.  At the package level, a minimum
severity can be set.  Messages less severe than this minimum logging level will
not be output.

Color Expressions

In addition to the standard printf-style formatting options (as defined in the
standard fmt package), this package supports inline expressions that control the
output of ANSI terminal escape sequences.  These expressions allow for a simple
mechanism to colorize log output, as well as applying graphical effects like
bold, underline, and blinking text (for terminals that support it).

By default, color expressions will only be honored if os.Stdin is attached to a
pseudoterminal.  This is the case when the program is run on the command line
and is not piped or redirected to another file.  This default ensures that the
colors are visible only in a visual context, but do not corrupt files or
pipelines with ANSI escape sequences.  Color sequences can be explicitly enabled
or disabled by setting the EnableColorExpressions package variable.

Using color expressions in format strings is done by wrapping the expression in
${expr}.  The general format for color expressions is:

	foregroundColor[+attributes[:backgroundColor[+attributes]]]

Colors (foreground and background):
	black
	red
	green
	yellow
	blue
	magenta
	cyan
	white
	[0-255]: numeric 8-bit color (for 256 color terminals)
	reset: Reset all color and graphics attributes to their defaults

Foreground Attributes:
	b: bold text
	B: blinking text
	h: high-intensity (bright text)
	i: inverted/reverse colors
	s: strikethrough
	u: underline

Background Attributes:
	h: high-intensity (bright text)

Examples

	log.Info("Hello, world!")
	// [11:22:33 0001] INFO Hello, world!

	log.Warningf("The %q operation could not be completed.", "add")
	// [11:22:33 0002] WARN The "add" operation could not be completed.

	log.Errorf("There was an ${red}error${reset} opening file ${blue+b:white}%s${reset}", filename)
	// [11:22:33 0003] ERRO There was an error opening file /tmp/file.txt
	//                                   ^^^^^              ^^^^^^^^^^^^^
	//                                   red text           blue text on white background
*/
package log

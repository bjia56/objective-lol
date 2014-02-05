import java.util.ArrayList;
import java.util.List;

import org.objectivelol.lang.LOLError;
import org.objectivelol.vm.RuntimeEnvironment;

public class MainClass {

	/**
	 * Command line invocation in progress
	 * TODO
	 * 
	 * @param args
	 * @throws LOLError 
	 */
	public static void main(String[] args) throws LOLError {
		Character c;
		Getopt.LongOption[] longopts = new Getopt.LongOption[] {
				new Getopt.LongOption("help", false, 'h')
		};
		
		List<String> sources = new ArrayList<String>();

		while((c = Getopt.getopt(args, "h", longopts)) != null) {
			switch(c) {
			case 'h':
				// TODO: help
			case ':':
				System.out.println("Error: Parameter required for " + args[Getopt.getIndex()] + "\nUse -h or --help for more information about options and required parameters.");
				System.exit(1);
			case '?':
				sources.add(args[Getopt.getIndex()]);
				break;
			default:
				System.out.println("Error: Internal error while parsing command line arguments. Parser returned: " + c);
				System.exit(1);
			}
		}
		
		longopts = null;
		c = null;
		
		RuntimeEnvironment re = RuntimeEnvironment.getRuntime();
		re.loadSource(sources.toArray(new String[0]));
		
		re.execute();
	}

	private static class Getopt {

		public static class LongOption {

			private String optionName;
			private boolean requiresArgument;
			private char code;

			public LongOption(String optionName, boolean requiresArgument, char code) {
				this.optionName = optionName;
				this.requiresArgument = requiresArgument;
				this.code = code;
			}

		}

		private static int index = -1;
		private static String param = null;
		
		public static int getIndex() {
			return index;
		}
		
		public static String getParam() {
			return param;
		}

		public static void reset() {
			index = -1;
		}

		public static Character getopt(final String[] args, final String shortopts, final LongOption[] longopts) {
			param = null;
			++index;

			if(index >= args.length || args[index].equals("")) {
				return null;
			}

			if(args[index].charAt(0) == '-') {
				if(args[index].length() > 1 && args[index].charAt(1) == '-') {
					for(LongOption lo : longopts) {
						if(args[index].substring(2).equals(lo.optionName)) {
							if(lo.requiresArgument) {
								if(index < args.length - 1) {
									++index;
									param = args[index + 1];
									return lo.code;
								} else {
									return ':';
								}
							} else {
								return lo.code;
							}
						}
					}
				}
				
				for(int i = 0; i < shortopts.length(); ++i) {
					if(args[index].substring(1).equals(shortopts.charAt(i) + "")) {
						if(i < shortopts.length() - 1 && shortopts.charAt(i + 1) == ':') {
							if(index < args.length - 1) {
								++index;
								param = args[index + 1];
								return shortopts.charAt(i);
							} else {
								return ':';
							}
						} else {
							return shortopts.charAt(i);
						}
					}
				}
			}

			return '?';
		}

	}

}

import java.io.File;
import java.io.IOException;
import java.util.ArrayList;
import java.util.List;

import org.objectivelol.lang.LOLError;
import org.objectivelol.vm.RuntimeEnvironment;

public class MainClass {
	
	private static final String version = "0.9.0";

	/**
	 * Command line invocation in progress
	 * TODO
	 * 
	 * @param args
	 * @throws LOLError 
	 * @throws IOException 
	 */
	public static void main(String[] args) throws LOLError, IOException {
		Character c;
		Getopt.LongOption[] longopts = new Getopt.LongOption[] {
				new Getopt.LongOption("help", false, 'h'),
				new Getopt.LongOption("version", false, 'v'),
				new Getopt.LongOption("lib", true, 'l'),
				new Getopt.LongOption("dir", true, 'd')
		};

		RuntimeEnvironment re = null;
		List<File> sources = new ArrayList<File>();
		String execDir = null;

		while((c = Getopt.getopt(args, "hvl:d:", longopts)) != null) {
			switch(c) {
			case 'h':
				// TODO: help
				System.exit(0);
			case 'v': // prints version information
				System.out.println("Objective-LOL Virtual Machine, version " + version);
				System.exit(0);
			case 'l': // sets the library directory
				re = RuntimeEnvironment.getRuntime(new File(Getopt.getParam()));
				break;
			case 'd': // sets the runtime directory
				execDir = Getopt.getParam();
				break;
			case ':': // parameter required but not found
				System.err.println("Error: Parameter required for " + args[Getopt.getIndex()] + "\nUse -h or --help for more information about options and required parameters.");
				System.exit(1);
			case '?': // unrecognized argument, assumed to be an input file
				String tmp = args[Getopt.getIndex()];
				File currentDirectory;

				// if the argument contains an explicit path, look for it
				if(tmp.contains(File.separator)) {
					// handle home directory
					if(tmp.startsWith("~")) {
						tmp = System.getProperty("user.home") + tmp.substring(1);
					}

					currentDirectory = new File(tmp).getParentFile();
					
					if(currentDirectory == null) {
						System.err.println("Error: Objective-LOL file " + tmp + " is invalid");
						System.exit(1);
					}
					
					tmp = new File(tmp).getName();
				} else {
					currentDirectory = new File(System.getProperty("user.dir"));
				}

				if(tmp.contains("*")) {
					tmp = tmp.replaceAll("\\.", "\\\\.").replaceAll("\\*", ".*");

					for(File f : currentDirectory.listFiles()) {
						if(!f.isFile()) {
							continue;
						}

						if(f.getName().matches(tmp)) {
							sources.add(f);
						}
					}
				} else {
					sources.add(new File(tmp));
				}

				break;
			default:
				System.err.println("Error: Internal error while parsing command line arguments. Parser returned: " + c);
				System.exit(1);
			}
		}

		if(re == null) {
			re = RuntimeEnvironment.getRuntime();
		}

		if(execDir != null) {
			re.setExecDir(new File(execDir));
		}


		c = null;
		longopts = null;
		execDir = null;
		
		re.loadSource(sources.toArray(new File[sources.size()]));
		
		sources = null;
		
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
									param = args[++index];
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
								param = args[++index];
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

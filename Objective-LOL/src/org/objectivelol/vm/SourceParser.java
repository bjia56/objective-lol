package org.objectivelol.vm;

import java.io.BufferedReader;
import java.io.File;
import java.io.FileNotFoundException;
import java.io.FileReader;
import java.io.IOException;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.LinkedHashMap;

import org.objectivelol.lang.LOLClass;
import org.objectivelol.lang.LOLError;
import org.objectivelol.lang.LOLFunction;
import org.objectivelol.lang.LOLNothing;
import org.objectivelol.lang.LOLObject;
import org.objectivelol.lang.LOLSource;
import org.objectivelol.lang.LOLValue;

public class SourceParser {

	private BufferedReader reader;
	private String fileName;

	private static String pointlessChar = (char)600 + "";

	public SourceParser(File file) {
		if(!file.isFile() || !file.getName().toLowerCase().endsWith(".lol")) {
			throw new IllegalArgumentException("Input file is not an Objective-LOL source file");
		}

		try {
			reader = new BufferedReader(new FileReader(file));
			fileName = file.getName().substring(0, file.getName().length() - 4).toUpperCase();
		} catch(FileNotFoundException e) {
			throw new RuntimeException("An unexpected IO error has occurred");
		}
	}

	public LOLSource parse() throws LOLError {
		try {
			String line;
			int lineNumber = 0;

			HashMap<String, ValueStruct> globalVariables = new HashMap<String, ValueStruct>();
			HashMap<String, LOLFunction> globalFunctions = new HashMap<String, LOLFunction>();
			HashMap<String, LOLClass> globalClasses = new HashMap<String, LOLClass>();

			while((line = reader.readLine()) != null) {
				lineNumber++;
				line = line.replaceAll("\\s+", " ").trim().toUpperCase();

				if(line.equals("")) {
					continue;
				}

				if(line.startsWith("HAI ME TEH VARIABLE") || line.startsWith("HAI ME TEH LOCKD VARIABLE")) {
					String[] tokens = line.split(" ");

					String type;
					LOLValue value;
					boolean isLocked;
					int offset = 0;

					if(isLocked = tokens[3].equals("LOCKD")) {
						offset = 1;
					}

					if(tokens.length == 4 + offset) {
						throw new LOLError("Line " + lineNumber + ": Variable identifier expected");
					}

					if(tokens.length < 7 + offset) {
						throw new LOLError("Line " + lineNumber + ": Variable type expected");
					}

					if(globalVariables.containsKey(tokens[4 + offset])) {
						throw new LOLError("Line " + lineNumber + ": Duplicate variable identifier detected");
					}

					type = tokens[6 + offset];

					if(tokens.length == 7 + offset) {
						value = LOLNothing.NOTHIN;
					} else {
						if(!tokens[7 + offset].equals("ITZ")) {
							throw new LOLError("Line "  + lineNumber + ": Unexpected symbol detected");
						}

						if(tokens.length < 9 + offset) {
							throw new LOLError("Line " + lineNumber + ": Value to assign expected");
						}

						if(tokens[8 + offset].equals("NEW")) {
							if(tokens.length < 10 + offset) {
								throw new LOLError("Line " + lineNumber + ": New object type expected");
							}

							if(!tokens[9 + offset].equals(type)) {
								throw new LOLError("Line " + lineNumber + ": Cannot instantiate specified object type into specified variable type");
							}

							if(tokens.length > 10 + offset) {
								if(tokens[10 + offset].equals("IN")) {
									if(tokens.length < 12 + offset) {
										throw new LOLError("Line " + lineNumber + ": Source expected at end of line");
									}

									if(tokens.length > 12 + offset) {
										throw new LOLError("Line " + lineNumber + ": Invalid symbols detected after new object type");
									}

									value = new LOLObjectRuntimeWrapper(tokens[11 + offset], type);
								}

								throw new LOLError("Line " + lineNumber + ": Invalid symbols detected after new object type");
							}

							value = new LOLObjectRuntimeWrapper(fileName, type);
						} else {
							StringBuilder s = new StringBuilder();

							for(int i = offset; i + 8 < tokens.length; i++) {
								if(i != offset) {
									s.append(" ");
								}
								s.append(tokens[8 + i]);
							}

							value = LOLValue.valueOf(s.toString().replace("\\\"", pointlessChar).replace("\"", "").replace(pointlessChar, "\"")).cast(tokens[6 + offset]);
						}
					}

					globalVariables.put(tokens[4 + offset], new ValueStruct(type, value, isLocked));

					continue;
				}

				if(line.startsWith("HAI ME TEH FUNCSHUN") || line.startsWith("HAI ME TEH NATIV FUNCSHUN")) {
					String[] tokens = line.split(" ");

					int offset = 0;

					if(tokens[3].equals("NATIV")) {
						offset = 1;
					}

					if(tokens.length == 4 + offset) {
						throw new LOLError("Line " + lineNumber + ": Function identifier expected");
					}

					if(globalFunctions.containsKey(tokens[4 + offset])) {
						throw new LOLError("Line " + lineNumber + ": Duplicate function identifier detected");
					}

					LinkedHashMap<String, String> fArgs = new LinkedHashMap<String, String>();
					String returnType = null;
					StringBuilder fCode = new StringBuilder();

					if(tokens.length > 5 + offset) {
						if(tokens[5 + offset].equals("TEH")) {
							if(tokens.length == 6 + offset) {
								throw new LOLError("Line " + lineNumber + ": Function return type expected");
							}

							returnType = tokens[6 + offset];

							if(tokens.length > 7 + offset) {
								if(tokens[7 + offset].equals("WIT")) {
									if(tokens.length == 8 + offset) {
										throw new LOLError("Line " + lineNumber + ": Argument identifier expected");
									}

									int index = 8 + offset;

									while(index < tokens.length) {
										if(index + 2 >= tokens.length) {
											throw new LOLError("Line " + lineNumber + ": Argument type expected");
										}

										if(fArgs.containsKey(tokens[index])) {
											throw new LOLError("Line " + lineNumber + ": Duplicate argument identifier detected");
										}

										fArgs.put(tokens[index], tokens[index + 2]);
										index += 5;
									}
								} else {
									throw new LOLError("Line " + lineNumber + ": Unexpected symbol detected");
								}
							}
						} else if(tokens[5 + offset].equals("WIT")) {
							if(tokens.length == 6 + offset) {
								throw new LOLError("Line " + lineNumber + ": Argument identifier expected");
							}

							int index = 6 + offset;

							while(index < tokens.length) {
								if(index + 2 >= tokens.length) {
									throw new LOLError("Line " + lineNumber + ": Argument type expected");
								}

								if(fArgs.containsKey(tokens[index])) {
									throw new LOLError("Line " + lineNumber + ": Duplicate argument identifier detected");
								}

								fArgs.put(tokens[index], tokens[index + 2]);
								index += 5;
							}
						} else {
							throw new LOLError("Line " + lineNumber + ": Unexpected symbol detected");
						}
					}

					if(offset == 0) {
						boolean first = true;
						while((line = reader.readLine()) != null && !line.startsWith("KTHXBAI")) {
							lineNumber++;
							line = line.replaceAll("\\s+", " ").trim().toUpperCase();

							if(line.equals("")) {
								continue;
							}

							fCode.append((first ? "" : "\n") + line);
							first = false;
						}

						if(line == null) {
							throw new LOLError("Unexpected EOF reached");
						}

						lineNumber++;

						if(!line.equals("KTHXBAI")) {
							throw new LOLError("Line " + lineNumber + ": Unexpected symbol detected");
						}

						globalFunctions.put(tokens[4], new LOLFunction(tokens[4], returnType, fArgs, null, null, fileName, fCode.toString()));
					} else {
						globalFunctions.put(tokens[5], new LOLFunction(tokens[5], returnType, fArgs, null, null, fileName, null) {

							@Override
							protected LOLValue run(LOLObject owner, LinkedHashMap<String, ValueStruct> args) throws LOLError {
								ArrayList<LOLValue> arguments = new ArrayList<LOLValue>();

								for(ValueStruct vs : args.values()) {
									arguments.add(vs.getValue());
								}

								return RuntimeEnvironment.getRuntime().getNative(getParentSource()).invoke(getName(), arguments.toArray(new LOLValue[0]));
							}

						});
					}

					continue;
				}

				if(line.startsWith("HAI ME TEH CLAS")) {
					String[] tokens = line.split(" ");

					if(tokens.length == 4) {
						throw new LOLError("Line " + lineNumber + ": Class identifier expected");
					}

					if(globalClasses.containsKey(tokens[4])) {
						throw new LOLError("Line " + lineNumber + ": Duplicate class identifier detected");
					}

					HashMap<String, ValueStruct> publicMemberVariables = new HashMap<String, ValueStruct>();
					HashMap<String, ValueStruct> privateMemberVariables = new HashMap<String, ValueStruct>();
					HashMap<String, ValueStruct> publicSharedVariables = new HashMap<String, ValueStruct>();
					HashMap<String, ValueStruct> privateSharedVariables = new HashMap<String, ValueStruct>();
					HashMap<String, LOLFunction> publicMemberFunctions = new HashMap<String, LOLFunction>();
					HashMap<String, LOLFunction> privateMemberFunctions = new HashMap<String, LOLFunction>();
					HashMap<String, LOLFunction> publicSharedFunctions = new HashMap<String, LOLFunction>();
					HashMap<String, LOLFunction> privateSharedFunctions = new HashMap<String, LOLFunction>();

					boolean isPublic = true;

					while((line = reader.readLine()) != null && !line.startsWith("KTHXBAI")) {
						lineNumber++;
						line = line.replaceAll("\\s+", " ").trim().toUpperCase();

						if(line.equals("")) {
							continue;
						}

						if(line.equals("EVRYONE")) {
							isPublic = true;
							continue;
						}

						if(line.equals("MAHSELF")) {
							isPublic = false;
							continue;
						}

						if(line.startsWith("DIS TEH VARIABLE") || line.startsWith("DIS TEH LOCKD VARIABLE") || line.startsWith("DIS TEH LOCKD SHARD VARIABLE") || line.startsWith("DIS TEH SHARD LOCKD VARIABLE")) {
							String[] tokens1 = line.split(" ");

							String type;
							LOLValue value;
							boolean isLocked;
							boolean isShared;
							int offset = 0;

							if(isLocked = (tokens1[2].equals("LOCKD") || tokens1[3].equals("LOCKD"))) {
								offset++;
							}

							if(isShared = (tokens1[2].equals("SHARD") || tokens1[3].equals("SHARD"))) {
								offset++;
							}

							if(tokens1[2].equals(tokens1[3])) {
								throw new LOLError("Line " + lineNumber + ": Unexpected symbol detected");
							}

							if(tokens1.length == 3 + offset) {
								throw new LOLError("Line " + lineNumber + ": Variable identifier expected");
							}

							if(tokens1.length < 6 + offset) {
								throw new LOLError("Line " + lineNumber + ": Variable type expected");
							}

							if(publicMemberVariables.containsKey(tokens1[3 + offset]) || privateMemberVariables.containsKey(tokens1[3 + offset]) || publicSharedVariables.containsKey(tokens1[3 + offset]) || privateSharedVariables.containsKey(tokens1[3 + offset])) {
								throw new LOLError("Line " + lineNumber + ": Duplicate variable identifier detected");
							}

							type = tokens1[5 + offset];

							if(tokens1.length == 6 + offset) {
								value = LOLNothing.NOTHIN;
							} else {
								if(!tokens1[6 + offset].equals("ITZ")) {
									throw new LOLError("Line "  + lineNumber + ": Unexpected symbol detected");
								}

								if(tokens1.length < 8 + offset) {
									throw new LOLError("Line " + lineNumber + ": Value to assign expected");
								}

								if(tokens1[7 + offset].equals("NEW")) {
									if(tokens1.length < 9 + offset) {
										throw new LOLError("Line " + lineNumber + ": New object type expected");
									}

									if(!tokens1[8 + offset].equals(type)) {
										throw new LOLError("Line " + lineNumber + ": Cannot instantiate specified object type into specified variable type");
									}

									if(tokens1.length > 9 + offset) {
										if(tokens1[9 + offset].equals("IN")) {
											if(tokens1.length < 11 + offset) {
												throw new LOLError("Line " + lineNumber + ": Source expected at end of line");
											}

											if(tokens1.length > 11 + offset) {
												throw new LOLError("Line " + lineNumber + ": Invalid symbols detected after new object type");
											}

											value = new LOLObjectRuntimeWrapper(tokens1[10 + offset], type);
										}

										throw new LOLError("Line " + lineNumber + ": Invalid symbols detected after new object type");
									}

									value = new LOLObjectRuntimeWrapper(fileName, type);
								} else {
									StringBuilder s = new StringBuilder();

									for(int i = offset; i + 7 < tokens1.length; i++) {
										if(i != offset) {
											s.append(" ");
										}
										s.append(tokens1[7 + i]);
									}

									value = LOLValue.valueOf(s.toString().replace("\\\"", pointlessChar).replace("\"", "").replace(pointlessChar, "\"")).cast(tokens1[5 + offset]);
								}
							}

							ValueStruct struct = new ValueStruct(type, value, isLocked);

							if(isPublic) {
								if(isShared) {
									publicSharedVariables.put(tokens1[3 + offset], struct);
								} else {
									publicMemberVariables.put(tokens1[3 + offset], struct);
								}
							} else {
								if(isShared) {
									privateSharedVariables.put(tokens1[3 + offset], struct);
								} else {
									privateMemberVariables.put(tokens1[3 + offset], struct);
								}
							}

							continue;
						}

						if(line.startsWith("DIS TEH FUNCSHUN") || line.startsWith("DIS TEH SHARD FUNCSHUN")) {
							String[] tokens1 = line.split(" ");

							boolean isShared;
							int offset = 0;

							if(isShared = tokens1[2].equals("SHARD")) {
								offset = 1;
							}

							if(tokens1.length == 3 + offset) {
								throw new LOLError("Line " + lineNumber + ": Function identifier expected");
							}

							if(publicMemberFunctions.containsKey(tokens1[3 + offset]) || privateMemberFunctions.containsKey(tokens1[3 + offset]) || publicSharedFunctions.containsKey(tokens1[3 + offset]) || privateSharedFunctions.containsKey(tokens1[3 + offset])) {
								throw new LOLError("Line " + lineNumber + ": Duplicate function identifier detected");
							}

							LinkedHashMap<String, String> fArgs = new LinkedHashMap<String, String>();
							String returnType = null;
							StringBuilder fCode = new StringBuilder();

							if(tokens1.length > 4 + offset) {
								if(tokens1[4 + offset].equals("TEH")) {
									if(tokens1.length == 5 + offset) {
										throw new LOLError("Line " + lineNumber + ": Function return type expected");
									}

									returnType = tokens1[5 + offset];

									if(tokens1.length > 6 + offset) {
										if(tokens1[6 + offset].equals("WIT")) {
											if(tokens1.length == 7) {
												throw new LOLError("Line " + lineNumber + ": Argument identifier expected");
											}

											int index = 7 + offset;

											while(index < tokens1.length) {
												if(index + 2 >= tokens1.length) {
													throw new LOLError("Line " + lineNumber + ": Argument type expected");
												}

												if(fArgs.containsKey(tokens1[index])) {
													throw new LOLError("Line " + lineNumber + ": Duplicate argument identifier detected");
												}

												fArgs.put(tokens1[index], tokens1[index + 2]);
												index += 5;
											}
										} else {
											throw new LOLError("Line " + lineNumber + ": Unexpected symbol detected");
										}
									}
								} else if(tokens1[4 + offset].equals("WIT")) {
									if(tokens1.length == 5 + offset) {
										throw new LOLError("Line " + lineNumber + ": Argument identifier expected");
									}

									int index = 5 + offset;

									while(index < tokens1.length) {
										if(index + 2 >= tokens1.length) {
											throw new LOLError("Line " + lineNumber + ": Argument type expected");
										}

										if(fArgs.containsKey(tokens1[index])) {
											throw new LOLError("Line " + lineNumber + ": Duplicate argument identifier detected");
										}

										fArgs.put(tokens1[index], tokens1[index + 2]);
										index += 5;
									}
								} else {
									throw new LOLError("Line " + lineNumber + ": Unexpected symbol detected");
								}
							}

							int nests = 1;

							boolean first = true;
							while((line = reader.readLine()) != null) {
								lineNumber++;
								line = line.replaceAll("\\s+", " ").trim().toUpperCase();

								if(line.equals("")) {
									continue;
								}

								if(line.startsWith("KTHX")) {
									if(!line.equals("KTHX")) {
										throw new LOLError("Line " + lineNumber + ": Unexpected symbol detected");
									}

									nests--;
								}

								if(line.startsWith("IZ") || line.startsWith("WHILE")) {
									nests++;
								}

								if(nests == 0) {
									break;
								}

								fCode.append((first ? "" : "\n") + line);
								first = false;
							}

							if(line == null) {
								throw new LOLError("Unexpected EOF reached");
							}

							if(isPublic) {
								if(isShared) {
									publicSharedFunctions.put(tokens1[3 + offset], new LOLFunction(tokens1[3 + offset], returnType, fArgs, true, tokens[4], fileName, fCode.toString()));
								} else {
									publicMemberFunctions.put(tokens1[3 + offset], new LOLFunction(tokens1[3 + offset], returnType, fArgs, false, tokens[4], fileName, fCode.toString()));
								}
							} else {
								if(isShared) {
									privateSharedFunctions.put(tokens1[3 + offset], new LOLFunction(tokens1[3 + offset], returnType, fArgs, true, tokens[4], fileName, fCode.toString()));
								} else {
									privateMemberFunctions.put(tokens1[3 + offset], new LOLFunction(tokens1[3 + offset], returnType, fArgs, false, tokens[4], fileName, fCode.toString()));
								}
							}

						}
					}

					globalClasses.put(tokens[4], new LOLClass(tokens[4], publicMemberVariables, privateMemberVariables, publicSharedVariables, privateSharedVariables, publicMemberFunctions, privateMemberFunctions, publicSharedFunctions, privateSharedFunctions, fileName));

					continue;
				}

				throw new LOLError("Line " + lineNumber + ": Unexpected symbol detected");
			}

			return new LOLSource(fileName, globalVariables, globalFunctions, globalClasses);
		} catch(IOException e) {
			throw new RuntimeException("An IO error has occurred");
		}
	}

}
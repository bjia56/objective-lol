package org.objectivelol.vm;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.StringReader;
import java.util.ArrayList;
import java.util.List;

import org.objectivelol.lang.LOLError;
import org.objectivelol.lang.LOLFunction;
import org.objectivelol.lang.LOLString;
import org.objectivelol.lang.LOLValue;

public class Parser {

	private static String pointlessChar = (char)600 + "";

	public static Expression parse(String s, LOLFunction context) throws LOLError {
		BufferedReader br = new BufferedReader(new StringReader(s));

		try {
			return parseBlock(br, context);
		} catch(IOException e) {
			throw new LOLError("An unexpected IO error has occurred");
		}
	}

	private static Expression parseBlock(BufferedReader br, LOLFunction context) throws IOException, LOLError {
		String line;

		ArrayList<Expression> statements = new ArrayList<Expression>();

		while((line = br.readLine()) != null) {
			line = line.replaceAll("\\s+", " ").trim();

			if(line.equals("")) {
				continue;
			}

			if(line.startsWith("BTW")) {
				continue;
			}

			if(line.contains(" BTW")) {
				line = line.substring(0, line.indexOf(" BTW"));
			}

			if(line.startsWith("IZ")) {
				if(!line.contains("?")) {
					throw new LOLError("Condition of IZ statement must be terminated by '?'");
				}

				if(line.indexOf("?") != line.length() - 1) {
					throw new LOLError("Unexpected symbol detected");
				}

				line = line.substring(2, line.length() - 1).trim();

				StringBuilder trueOperations = new StringBuilder();
				StringBuilder elseOperations = new StringBuilder();

				boolean first = true, isElse = false;
				String line2;
				while(!(line2 = br.readLine()).startsWith("KTHX")) {
					if(line2.trim().equals("")) {
						continue;
					}

					if(line2.startsWith("NOPE")) {
						if(!line2.equals("NOPE")) {
							throw new LOLError("Unexpected symbol detected");
						}

						isElse = true;
						first = true;

						continue;
					}

					if(!isElse) {
						trueOperations.append((first ? "" : "\n") + line2);
					} else {
						elseOperations.append((first ? "" : "\n") + line2);
					}

					first = false;
				}

				if(!line2.equals("KTHX")) {
					throw new LOLError("Unexpected symbol detected");
				}

				Expression condition = parseLine(line, context);
				Expression code = parseBlock(new BufferedReader(new StringReader(trueOperations.toString())), context);
				Expression code2 = (isElse ? parseBlock(new BufferedReader(new StringReader(elseOperations.toString())), context) : null);

				statements.add(new IfStatement(condition, code, code2));
				continue;
			}

			if(line.startsWith("WHILE")) {
				line = line.substring(5).trim();

				StringBuilder sb = new StringBuilder();

				boolean first = true;
				String line2;
				while(!(line2 = br.readLine()).startsWith("KTHX")) {
					if(line2.trim().equals("")) {
						continue;
					}

					sb.append((first ? "" : "\n") + line2);
					first = false;
				}

				if(!line2.equals("KTHX")) {
					throw new LOLError("Unexpected symbol detected");
				}

				Expression condition = parseLine(line, context);
				Expression code = parseBlock(new BufferedReader(new StringReader(sb.toString())), context);

				statements.add(new WhileStatement(condition, code));
				continue;
			}

			statements.add(parseLine(line, context));
		}

		return new StatementBlock(statements);
	}

	private static Expression parseLine(String line, LOLFunction context) throws LOLError {
		String[] tmp = line.split(" ");
		List<String> tokens = new ArrayList<String>();

		boolean isStringLiteral = false;
		boolean firstPartOfLiteral = false;
		StringBuilder literal = null;
		for(String s : tmp) {
			if(firstPartOfLiteral) {
				firstPartOfLiteral = false;
			}

			if(!isStringLiteral && s.charAt(0) == '\"') {
				isStringLiteral = true;
				firstPartOfLiteral = true;
				literal = new StringBuilder(s);
			}

			if(!isStringLiteral) {
				tokens.add(s);
				continue;
			} else {
				if(!firstPartOfLiteral) {
					literal.append(" " + s);
				}
			}

			if(isStringLiteral && s.charAt(s.length() - 1) == '\"' && (s.length() > 1 ? s.charAt(s.length() - 2) != '\\' : true)) {
				isStringLiteral = false;

				String tmpLiteral;
				if(firstPartOfLiteral) {
					tmpLiteral = literal.toString();
					firstPartOfLiteral = false;
				} else {
					tmpLiteral = literal.append(s).toString();
				}

				tokens.add('\"' + tmpLiteral.replace("\\\"", pointlessChar).replace("\"", "").replace(pointlessChar, "\"") + '\"');
				continue;
			}
		}

		if(isStringLiteral) {
			throw new LOLError("Termination character of string literal not found");
		}

		Expression argFunctionCall = null;

		if(tokens.contains("WIT")) {
			int functionStart = tokens.indexOf("WIT") - 1;

			if(functionStart < 0) {
				throw new LOLError("Function identifier expected before WIT");
			}

			if(tokens.size() - functionStart < 2) {
				throw new LOLError("No function arguments specified");
			}

			ArrayList<Expression> arguments = new ArrayList<Expression>();

			int index = functionStart + 2;
			while(index < tokens.size()) {
				StringBuilder argumentExpression = new StringBuilder();

				int i = index;
				for(; i < tokens.size(); i++) {
					if(tokens.get(i).equals("AN") && i + 1 < tokens.size() && tokens.get(i + 1).equals("WIT")) {
						break;
					}

					if(tokens.get(i).equals("WIT")) {
						throw new LOLError("Only one function call per line allowed");
					}

					argumentExpression.append((i == index ? "" : " ") + tokens.get(i));
				}

				if(tokens.size() - i == 2) {
					throw new LOLError("Function argument expected after AN WIT");
				}

				arguments.add(parseStatement(argumentExpression.toString(), null));

				index = i + 2;
			}

			if(functionStart >= 1) {
				if(tokens.get(functionStart - 1).equals("IN")) {
					if(functionStart - 1 == 0) {
						throw new LOLError("Function identifier expected before IN");
					}

					argFunctionCall = new MemberArgFunction(tokens.get(functionStart), tokens.get(functionStart - 2), arguments);
				}
			}

			if(argFunctionCall == null) {
				argFunctionCall = new ArgFunction(tokens.get(functionStart), arguments);
			}
		}

		if(tokens.get(0).equals("GIVEZ")) {
			if(tokens.size() == 1) {
				throw new LOLError("Expected expression after GIVEZ");
			}
			
			if(tokens.size() == 2 && tokens.get(1).equals("UP")) {
				return new Expression.Return(null);
			}
			
			if(tokens.contains("ITZ")) {
				throw new LOLError("GIVEZ line cannot include an assignment");
			}

			if(argFunctionCall == null) {
				StringBuilder expression = new StringBuilder();

				for(int i = 1; i < tokens.size(); i++) {
					expression.append((i == 1 ? "" : " ") + tokens.get(i));
				}

				return new Expression.Return(parseStatement(expression.toString(), null));
			} else {
				StringBuilder expression = new StringBuilder();

				for(int i = 1, functionStart = tokens.indexOf("WIT") - (tokens.contains("IN") ? 3 : 1); i < functionStart; i++) {
					expression.append((i == 1 ? "" : " ") + tokens.get(i));
				}

				return new Expression.Return(parseStatement(expression.toString(), argFunctionCall));
			}
		}

		if(line.startsWith("I HAS A") && tokens.get(2).equals("A")) {
			int lockedOffset = (tokens.contains("LOCKD") ? 1 : 0);

			if(tokens.size() < 4 + lockedOffset || !tokens.get(3 + lockedOffset).equals("VARIABLE")) {
				throw new LOLError("VARIABLE expected after I HAS A" + (lockedOffset == 1 ? " LOCKED" : ""));
			}

			if(tokens.size() < 5 + lockedOffset) {
				throw new LOLError("Variable identifier expected after VARIABLE");
			}

			if(tokens.size() < 6 + lockedOffset || !tokens.get(5 + lockedOffset).equals("TEH")) {
				throw new LOLError("TEH expected after variable identifier");
			}

			if(tokens.size() < 7 + lockedOffset) {
				throw new LOLError("Variable type expected after TEH");
			}

			if(tokens.contains("ITZ")) {
				if(tokens.indexOf("ITZ") != 7 + lockedOffset) {
					throw new LOLError("Incorrect positioning of ITZ in assignment to new variable");
				}

				if(tokens.size() == 8 + lockedOffset) {
					throw new LOLError("Expression to assign expected after ITZ");
				}

				if(tokens.get(8 + lockedOffset).equals("NEW")) {
					if(argFunctionCall != null) {
						throw new LOLError("Cannot have function call in line instantiating a new object");
					}

					if(tokens.size() < 10 + lockedOffset) {
						throw new LOLError("New object type expected");
					}

					if(!tokens.get(9 + lockedOffset).equals(tokens.get(6 + lockedOffset))) {
						throw new LOLError("Cannot instantiate specified object type into specified variable type");
					}

					if(tokens.size() > 10 + lockedOffset) {
						if(tokens.get(10 + lockedOffset).equals("IN")) {
							if(tokens.size() < 12 + lockedOffset) {
								throw new LOLError("Source expected at end of line");
							}

							if(tokens.size() > 12 + lockedOffset) {
								throw new LOLError("Invalid symbols detected after new object type");
							}

							return new DeclareVariable(tokens.get(4 + lockedOffset), tokens.get(6 + lockedOffset), lockedOffset == 1, new Value(new LOLObjectRuntimeWrapper(tokens.get(11 + lockedOffset), tokens.get(6 + lockedOffset))));
						}

						throw new LOLError("Invalid symbols detected after new object type");
					}

					return new DeclareVariable(tokens.get(4 + lockedOffset), tokens.get(6 + lockedOffset), lockedOffset == 1, new Value(new LOLObjectRuntimeWrapper(context.getParentSource(), tokens.get(6 + lockedOffset))));
				}

				if(argFunctionCall == null) {
					StringBuilder expression = new StringBuilder();

					for(int i = 8 + lockedOffset; i < tokens.size(); i++) {
						expression.append((i == 8 + lockedOffset ? "" : " ") + tokens.get(i));
					}

					return new DeclareVariable(tokens.get(4 + lockedOffset), tokens.get(6 + lockedOffset), lockedOffset == 1, parseStatement(expression.toString(), null));
				} else {
					StringBuilder expression = new StringBuilder();

					for(int i = 8 + lockedOffset, functionStart = tokens.indexOf("WIT") - (tokens.contains("IN") ? 3 : 1); i < functionStart; i++) {
						expression.append((i == 8 + lockedOffset ? "" : " ") + tokens.get(i));
					}

					return new DeclareVariable(tokens.get(4 + lockedOffset), tokens.get(6 + lockedOffset), lockedOffset == 1, parseStatement(expression.toString(), argFunctionCall));
				}
			} else {
				throw new LOLError("Declared variable must be initialized with a value");
			}
		}

		if(tokens.contains("ITZ")) {
			int target = tokens.indexOf("ITZ") - 1;

			if(target < 0) {
				throw new LOLError("Variable identifier expected before ITZ");
			}

			if(target != tokens.lastIndexOf("ITZ") - 1) {
				throw new LOLError("Only one assignment operator allowed per line");
			}

			if(target == 0) {
				if(tokens.size() < 3) {
					throw new LOLError("Expression to assign expected after ITZ");
				}

				if(tokens.get(2).equals("NEW")) {
					if(argFunctionCall != null) {
						throw new LOLError("Cannot have function call in line instantiating a new object");
					}

					if(tokens.size() < 4) {
						throw new LOLError("New object type expected");
					}

					if(tokens.size() > 4) {
						if(tokens.get(4).equals("IN")) {
							if(tokens.size() < 6) {
								throw new LOLError("Source expected at end of line");
							}

							if(tokens.size() > 6) {
								throw new LOLError("Invalid symbols detected after new object type");
							}

							return new SimpleAssignment(tokens.get(0), new Value(new LOLObjectRuntimeWrapper(tokens.get(5), tokens.get(3))));
						}

						throw new LOLError("Invalid symbols detected after new object type");
					}

					return new SimpleAssignment(tokens.get(0), new Value(new LOLObjectRuntimeWrapper(context.getParentSource(), tokens.get(3))));
				}

				if(argFunctionCall == null) {
					StringBuilder expression = new StringBuilder();

					for(int i = 2; i < tokens.size(); i++) {
						expression.append((i == 2 ? "" : " ") + tokens.get(i));
					}

					return new SimpleAssignment(tokens.get(0), parseStatement(expression.toString(), null));
				} else {
					StringBuilder expression = new StringBuilder();

					for(int i = 2, functionStart = tokens.indexOf("WIT") - (tokens.contains("IN") ? 3 : 1); i < functionStart; i++) {
						expression.append((i == 2 ? "" : " ") + tokens.get(i));
					}

					return new SimpleAssignment(tokens.get(0), parseStatement(expression.toString(), argFunctionCall));
				}
			}

			if(tokens.contains("IN") && tokens.indexOf("IN") < target) {
				int inLocation = tokens.indexOf("IN");

				if(inLocation == 0) {
					throw new LOLError("Variable identifier expected before IN");
				}

				if(inLocation > 1) {
					throw new LOLError("Unexpected tokens before IN in assignment line");
				}

				if(inLocation != target - 1) {
					throw new LOLError("Unexpected tokens between IN and ITZ in assignment line");
				}

				if(tokens.size() < 5) {
					throw new LOLError("Expression to assign expected after ITZ");
				}

				if(argFunctionCall == null) {
					StringBuilder expression = new StringBuilder();

					for(int i = 4; i < tokens.size(); i++) {
						expression.append((i == 4 ? "" : " ") + tokens.get(i));
					}

					return new ComplexAssignment(tokens.get(2), tokens.get(0), parseStatement(expression.toString(), null));
				} else {
					StringBuilder expression = new StringBuilder();

					for(int i = 4, functionStart = tokens.indexOf("WIT") - (tokens.contains("IN") ? 3 : 1); i < functionStart; i++) {
						expression.append((i == 4 ? "" : " ") + tokens.get(i));
					}

					return new ComplexAssignment(tokens.get(2), tokens.get(0), parseStatement(expression.toString(), argFunctionCall));
				}
			}
		}

		if(argFunctionCall == null) {
			return parseStatement(line, null);
		} else {
			StringBuilder expression = new StringBuilder();

			for(int i = 0, functionStart = tokens.indexOf("WIT") - (tokens.contains("IN") ? 3 : 1); i < functionStart; i++) {
				expression.append((i == 0 ? "" : " ") + tokens.get(i));
			}

			return parseStatement(expression.toString(), argFunctionCall);
		}
	}

	private static Expression parseStatement(String line, Expression function) throws LOLError {
		if(line.trim().equals("")) {
			return function;
		}

		String[] tmp = line.split(" ");
		List<Object> tokens = new ArrayList<Object>();

		boolean isStringLiteral = false;
		boolean firstPartOfLiteral = false;
		StringBuilder literal = null;
		for(String s : tmp) {
			if(firstPartOfLiteral) {
				firstPartOfLiteral = false;
			}

			if(!isStringLiteral && s.charAt(0) == '\"') {
				isStringLiteral = true;
				firstPartOfLiteral = true;
				literal = new StringBuilder(s);
			}

			if(!isStringLiteral) {
				tokens.add(s);
				continue;
			} else {
				if(!firstPartOfLiteral) {
					literal.append(" " + s);
				}
			}

			if(isStringLiteral && s.charAt(s.length() - 1) == '\"' && (s.length() > 1 ? s.charAt(s.length() - 2) != '\\' : true)) {
				isStringLiteral = false;

				String tmpLiteral;
				if(firstPartOfLiteral) {
					tmpLiteral = literal.toString();
					firstPartOfLiteral = false;
				} else {
					tmpLiteral = literal.append(s).toString();
				}

				tokens.add('\"' + tmpLiteral.replace("\\\"", pointlessChar).replace("\"", "").replace(pointlessChar, "\"") + '\"');
				continue;
			}
		}

		if(isStringLiteral) {
			throw new LOLError("Termination character of string literal not found");
		}

		if(function != null) {
			tokens.add(function);
		}

		while(tokens.contains("IN")) {
			int inIndex = tokens.indexOf("IN");

			if(inIndex == 0) {
				throw new LOLError("Member identifier expected before IN");
			}

			if(inIndex == tokens.size() - 1) {
				throw new LOLError("Source or class expected after IN");
			}

			if(!(tokens.get(inIndex + 1) instanceof String)) {
				throw new LOLError("Source or class after IN cannot be an expression");
			}

			if(!(tokens.get(inIndex - 1) instanceof String)) {
				throw new LOLError("Member before IN cannot be an expression");
			}

			tokens.add(inIndex - 1, new MemberVariableAndNoArgFunction((String)tokens.get(inIndex + 1), (String)tokens.get(inIndex - 1)));
			tokens.remove(inIndex);
			tokens.remove(inIndex);
			tokens.remove(inIndex);
		}

		while(tokens.contains("AS")) {
			int asLocation = tokens.indexOf("AS");

			if(asLocation == 0) {
				throw new LOLError("Variable to cast expected before AS");
			}

			if(tokens.indexOf("SAEM") == asLocation - 1) {
				tokens.set(asLocation, "AS_SAEM");
				continue;
			}

			if(asLocation == tokens.size() - 1) {
				throw new LOLError("Target cast type expected after AS");
			}

			if(!(tokens.get(asLocation + 1) instanceof String)) {
				throw new LOLError("Target cast type must not be a function");
			}

			if(!(tokens.get(asLocation - 1) instanceof String)) {
				throw new LOLError("Variable to cast must not be an expression");
			}

			Expression variable = null;
			String varString = (String)tokens.get(asLocation - 1);

			if(varString.charAt(0) == '\"') {
				variable = new Value(new LOLString(varString.substring(1, varString.length() - 1)));
			} else {
				try {
					Double.parseDouble(varString);
				} catch(NumberFormatException e) {
					try {
						Long.parseLong(varString);
					} catch(NumberFormatException e2) {
						try {
							//varString = varString.toUpperCase();
							if(!varString.startsWith("0X")) {
								throw new NumberFormatException();
							}
							Long.parseLong(varString.replaceFirst("0X", ""), 16);
						} catch(NumberFormatException e3) {
							if(!varString.equals("YEZ") && !varString.equals("NO")) {
								variable = new VariableAndNoArgFunction(varString);
							}
						}
					}
				}
			}

			if(variable == null) {
				variable = new Value(LOLValue.valueOf(varString));
			}

			tokens.add(asLocation - 1, new Cast(variable , (String)tokens.get(asLocation + 1)));
			tokens.remove(asLocation);
			tokens.remove(asLocation);
			tokens.remove(asLocation);
		}

		while(tokens.contains("TIEMZ") || tokens.contains("DIVIDEZ")) {
			int timesIndex = tokens.indexOf("TIEMZ");
			int dividesIndex = tokens.indexOf("DIVIDEZ");

			if(timesIndex == -1 || (dividesIndex < timesIndex && dividesIndex != -1)) {
				if(dividesIndex == 0 || dividesIndex == tokens.size() - 1) {
					throw new LOLError("Two arguments required for DIVIDEZ operation");
				}

				Expression expressionBefore = (tokens.get(dividesIndex - 1) instanceof Expression ? (Expression)tokens.get(dividesIndex - 1) : null);

				if(expressionBefore == null) {
					String expStringBefore = (String)tokens.get(dividesIndex - 1);

					if(expStringBefore.charAt(0) == '\"') {
						expressionBefore = new Value(new LOLString(expStringBefore.substring(1, expStringBefore.length() - 1)));
					} else {
						try {
							Double.parseDouble(expStringBefore);
						} catch(NumberFormatException e) {
							try {
								Long.parseLong(expStringBefore);
							} catch(NumberFormatException e2) {
								try {
									//expStringBefore = expStringBefore.toUpperCase();
									if(!expStringBefore.startsWith("0X")) {
										throw new NumberFormatException();
									}
									Long.parseLong(expStringBefore.replaceFirst("0X", ""), 16);
								} catch(NumberFormatException e3) {
									if(!expStringBefore.equals("YEZ") && !expStringBefore.equals("NO")) {
										expressionBefore = new VariableAndNoArgFunction(expStringBefore);
									}
								}
							}
						}
					}

					if(expressionBefore == null) {
						expressionBefore = new Value(LOLValue.valueOf(expStringBefore));
					}
				}

				Expression expressionAfter = (tokens.get(dividesIndex + 1) instanceof Expression ? (Expression)tokens.get(dividesIndex + 1) : null);

				if(expressionAfter == null) {
					String expStringAfter = (String)tokens.get(dividesIndex + 1);

					if(expStringAfter.charAt(0) == '\"') {
						expressionAfter = new Value(new LOLString(expStringAfter.substring(1, expStringAfter.length() - 1)));
					} else {
						try {
							Double.parseDouble(expStringAfter);
						} catch(NumberFormatException e) {
							try {
								Long.parseLong(expStringAfter);
							} catch(NumberFormatException e2) {
								try {
									//expStringAfter = expStringAfter.toUpperCase();
									if(!expStringAfter.startsWith("0X")) {
										throw new NumberFormatException();
									}
									Long.parseLong(expStringAfter.replaceFirst("0X", ""), 16);
								} catch(NumberFormatException e3) {
									if(!expStringAfter.equals("YEZ") && !expStringAfter.equals("NO")) {
										expressionAfter = new VariableAndNoArgFunction(expStringAfter);
									}
								}
							}
						}
					}

					if(expressionAfter == null) {
						expressionAfter = new Value(LOLValue.valueOf(expStringAfter));
					}
				}

				tokens.add(dividesIndex - 1, new Divide(expressionBefore, expressionAfter));
				tokens.remove(dividesIndex);
				tokens.remove(dividesIndex);
				tokens.remove(dividesIndex);
				continue;
			}

			if(dividesIndex == -1 || (timesIndex < dividesIndex && timesIndex != -1)) {
				if(timesIndex == 0 || timesIndex == tokens.size() - 1) {
					throw new LOLError("Two arguments required for TIEMZ operation");
				}

				Expression expressionBefore = (tokens.get(timesIndex - 1) instanceof Expression ? (Expression)tokens.get(timesIndex - 1) : null);

				if(expressionBefore == null) {
					String expStringBefore = (String)tokens.get(timesIndex - 1);

					if(expStringBefore.charAt(0) == '\"') {
						expressionBefore = new Value(new LOLString(expStringBefore.substring(1, expStringBefore.length() - 1)));
					} else {
						try {
							Double.parseDouble(expStringBefore);
						} catch(NumberFormatException e) {
							try {
								Long.parseLong(expStringBefore);
							} catch(NumberFormatException e2) {
								try {
									//expStringBefore = expStringBefore.toUpperCase();
									if(!expStringBefore.startsWith("0X")) {
										throw new NumberFormatException();
									}
									Long.parseLong(expStringBefore.replaceFirst("0X", ""), 16);
								} catch(NumberFormatException e3) {
									if(!expStringBefore.equals("YEZ") && !expStringBefore.equals("NO")) {
										expressionBefore = new VariableAndNoArgFunction(expStringBefore);
									}
								}
							}
						}
					}

					if(expressionBefore == null) {
						expressionBefore = new Value(LOLValue.valueOf(expStringBefore));
					}
				}

				Expression expressionAfter = (tokens.get(timesIndex + 1) instanceof Expression ? (Expression)tokens.get(timesIndex + 1) : null);

				if(expressionAfter == null) {
					String expStringAfter = (String)tokens.get(timesIndex + 1);

					if(expStringAfter.charAt(0) == '\"') {
						expressionAfter = new Value(new LOLString(expStringAfter.substring(1, expStringAfter.length() - 1)));
					} else {
						try {
							Double.parseDouble(expStringAfter);
						} catch(NumberFormatException e) {
							try {
								Long.parseLong(expStringAfter);
							} catch(NumberFormatException e2) {
								try {
									//expStringAfter = expStringAfter.toUpperCase();
									if(!expStringAfter.startsWith("0X")) {
										throw new NumberFormatException();
									}
									Long.parseLong(expStringAfter.replaceFirst("0X", ""), 16);
								} catch(NumberFormatException e3) {
									if(!expStringAfter.equals("YEZ") && !expStringAfter.equals("NO")) {
										expressionAfter = new VariableAndNoArgFunction(expStringAfter);
									}
								}
							}
						}
					}

					if(expressionAfter == null) {
						expressionAfter = new Value(LOLValue.valueOf(expStringAfter));
					}
				}

				tokens.add(timesIndex - 1, new Multiply(expressionBefore, expressionAfter));
				tokens.remove(timesIndex);
				tokens.remove(timesIndex);
				tokens.remove(timesIndex);
				continue;
			}
		}

		while(tokens.contains("MOAR") || tokens.contains("LES")) {
			int addIndex = tokens.indexOf("MOAR");
			int subtractIndex = tokens.indexOf("LES");

			if(addIndex == -1 || (subtractIndex < addIndex && subtractIndex != -1)) {
				if(subtractIndex == 0 || subtractIndex == tokens.size() - 1) {
					throw new LOLError("Two arguments required for LES operation");
				}

				Expression expressionBefore = (tokens.get(subtractIndex - 1) instanceof Expression ? (Expression)tokens.get(subtractIndex - 1) : null);

				if(expressionBefore == null) {
					String expStringBefore = (String)tokens.get(subtractIndex - 1);

					if(expStringBefore.charAt(0) == '\"') {
						expressionBefore = new Value(new LOLString(expStringBefore.substring(1, expStringBefore.length() - 1)));
					} else {
						try {
							Double.parseDouble(expStringBefore);
						} catch(NumberFormatException e) {
							try {
								Long.parseLong(expStringBefore);
							} catch(NumberFormatException e2) {
								try {
									//expStringBefore = expStringBefore.toUpperCase();
									if(!expStringBefore.startsWith("0X")) {
										throw new NumberFormatException();
									}
									Long.parseLong(expStringBefore.replaceFirst("0X", ""), 16);
								} catch(NumberFormatException e3) {
									if(!expStringBefore.equals("YEZ") && !expStringBefore.equals("NO")) {
										expressionBefore = new VariableAndNoArgFunction(expStringBefore);
									}
								}
							}
						}
					}

					if(expressionBefore == null) {
						expressionBefore = new Value(LOLValue.valueOf(expStringBefore));
					}
				}

				Expression expressionAfter = (tokens.get(subtractIndex + 1) instanceof Expression ? (Expression)tokens.get(subtractIndex + 1) : null);

				if(expressionAfter == null) {
					String expStringAfter = (String)tokens.get(subtractIndex + 1);

					if(expStringAfter.charAt(0) == '\"') {
						expressionAfter = new Value(new LOLString(expStringAfter.substring(1, expStringAfter.length() - 1)));
					} else {
						try {
							Double.parseDouble(expStringAfter);
						} catch(NumberFormatException e) {
							try {
								Long.parseLong(expStringAfter);
							} catch(NumberFormatException e2) {
								try {
									//expStringAfter = expStringAfter.toUpperCase();
									if(!expStringAfter.startsWith("0X")) {
										throw new NumberFormatException();
									}
									Long.parseLong(expStringAfter.replaceFirst("0X", ""), 16);
								} catch(NumberFormatException e3) {
									if(!expStringAfter.equals("YEZ") && !expStringAfter.equals("NO")) {
										expressionAfter = new VariableAndNoArgFunction(expStringAfter);
									}
								}
							}
						}
					}

					if(expressionAfter == null) {
						expressionAfter = new Value(LOLValue.valueOf(expStringAfter));
					}
				}

				tokens.add(subtractIndex - 1, new Subtract(expressionBefore, expressionAfter));
				tokens.remove(subtractIndex);
				tokens.remove(subtractIndex);
				tokens.remove(subtractIndex);
				continue;
			}

			if(subtractIndex == -1 || (addIndex < subtractIndex && addIndex != -1)) {
				if(addIndex == 0 || addIndex == tokens.size() - 1) {
					throw new LOLError("Two arguments required for MOAR operation");
				}

				Expression expressionBefore = (tokens.get(addIndex - 1) instanceof Expression ? (Expression)tokens.get(addIndex - 1) : null);

				if(expressionBefore == null) {
					String expStringBefore = (String)tokens.get(addIndex - 1);

					if(expStringBefore.charAt(0) == '\"') {
						expressionBefore = new Value(new LOLString(expStringBefore.substring(1, expStringBefore.length() - 1)));
					} else {
						try {
							Double.parseDouble(expStringBefore);
						} catch(NumberFormatException e) {
							try {
								Long.parseLong(expStringBefore);
							} catch(NumberFormatException e2) {
								try {
									//expStringBefore = expStringBefore.toUpperCase();
									if(!expStringBefore.startsWith("0X")) {
										throw new NumberFormatException();
									}
									Long.parseLong(expStringBefore.replaceFirst("0X", ""), 16);
								} catch(NumberFormatException e3) {
									if(!expStringBefore.equals("YEZ") && !expStringBefore.equals("NO")) {
										expressionBefore = new VariableAndNoArgFunction(expStringBefore);
									}
								}
							}
						}
					}

					if(expressionBefore == null) {
						expressionBefore = new Value(LOLValue.valueOf(expStringBefore));
					}
				}

				Expression expressionAfter = (tokens.get(addIndex + 1) instanceof Expression ? (Expression)tokens.get(addIndex + 1) : null);

				if(expressionAfter == null) {
					String expStringAfter = (String)tokens.get(addIndex + 1);

					if(expStringAfter.charAt(0) == '\"') {
						expressionAfter = new Value(new LOLString(expStringAfter.substring(1, expStringAfter.length() - 1)));
					} else {
						try {
							Double.parseDouble(expStringAfter);
						} catch(NumberFormatException e) {
							try {
								Long.parseLong(expStringAfter);
							} catch(NumberFormatException e2) {
								try {
									//expStringAfter = expStringAfter.toUpperCase();
									if(!expStringAfter.startsWith("0X")) {
										throw new NumberFormatException();
									}
									Long.parseLong(expStringAfter.replaceFirst("0X", ""), 16);
								} catch(NumberFormatException e3) {
									if(!expStringAfter.equals("YEZ") && !expStringAfter.equals("NO")) {
										expressionAfter = new VariableAndNoArgFunction(expStringAfter);
									}
								}
							}
						}
					}

					if(expressionAfter == null) {
						expressionAfter = new Value(LOLValue.valueOf(expStringAfter));
					}
				}

				tokens.add(addIndex - 1, new Add(expressionBefore, expressionAfter));
				tokens.remove(addIndex);
				tokens.remove(addIndex);
				tokens.remove(addIndex);
				continue;
			}
		}

		while(tokens.contains("BIGGR") || tokens.contains("SMALLR") || tokens.contains("THAN")) {
			int greaterThanIndex = tokens.indexOf("BIGGR");
			int lessThanIndex = tokens.indexOf("SMALLR");
			int thanIndex = tokens.indexOf("THAN");

			if(thanIndex == 0 || (greaterThanIndex == -1 && lessThanIndex == -1)) {
				throw new LOLError("BIGGR or SMALLR expected before THAN");
			}

			if(greaterThanIndex == -1 || (lessThanIndex < greaterThanIndex && lessThanIndex != -1)) {
				if(lessThanIndex + 1 != thanIndex) {
					throw new LOLError("THAN expected after SMALLR");
				}

				if(lessThanIndex == 0 || lessThanIndex == tokens.size() - 2) {
					throw new LOLError("Two arguments required for SMALLR THAN operation");
				}

				Expression expressionBefore = (tokens.get(lessThanIndex - 1) instanceof Expression ? (Expression)tokens.get(lessThanIndex - 1) : null);

				if(expressionBefore == null) {
					String expStringBefore = (String)tokens.get(lessThanIndex - 1);

					if(expStringBefore.charAt(0) == '\"') {
						expressionBefore = new Value(new LOLString(expStringBefore.substring(1, expStringBefore.length() - 1)));
					} else {
						try {
							Double.parseDouble(expStringBefore);
						} catch(NumberFormatException e) {
							try {
								Long.parseLong(expStringBefore);
							} catch(NumberFormatException e2) {
								try {
									//expStringBefore = expStringBefore.toUpperCase();
									if(!expStringBefore.startsWith("0X")) {
										throw new NumberFormatException();
									}
									Long.parseLong(expStringBefore.replaceFirst("0X", ""), 16);
								} catch(NumberFormatException e3) {
									if(!expStringBefore.equals("YEZ") && !expStringBefore.equals("NO")) {
										expressionBefore = new VariableAndNoArgFunction(expStringBefore);
									}
								}
							}
						}
					}

					if(expressionBefore == null) {
						expressionBefore = new Value(LOLValue.valueOf(expStringBefore));
					}
				}

				Expression expressionAfter = (tokens.get(lessThanIndex + 2) instanceof Expression ? (Expression)tokens.get(lessThanIndex + 2) : null);

				if(expressionAfter == null) {
					String expStringAfter = (String)tokens.get(lessThanIndex + 2);

					if(expStringAfter.charAt(0) == '\"') {
						expressionAfter = new Value(new LOLString(expStringAfter.substring(1, expStringAfter.length() - 1)));
					} else {
						try {
							Double.parseDouble(expStringAfter);
						} catch(NumberFormatException e) {
							try {
								Long.parseLong(expStringAfter);
							} catch(NumberFormatException e2) {
								try {
									//expStringAfter = expStringAfter.toUpperCase();
									if(!expStringAfter.startsWith("0X")) {
										throw new NumberFormatException();
									}
									Long.parseLong(expStringAfter.replaceFirst("0X", ""), 16);
								} catch(NumberFormatException e3) {
									if(!expStringAfter.equals("YEZ") && !expStringAfter.equals("NO")) {
										expressionAfter = new VariableAndNoArgFunction(expStringAfter);
									}
								}
							}
						}
					}

					if(expressionAfter == null) {
						expressionAfter = new Value(LOLValue.valueOf(expStringAfter));
					}
				}

				tokens.add(lessThanIndex - 1, new LessThan(expressionBefore, expressionAfter));
				tokens.remove(lessThanIndex);
				tokens.remove(lessThanIndex);
				tokens.remove(lessThanIndex);
				tokens.remove(lessThanIndex);
				continue;
			}

			if(lessThanIndex == -1 || (greaterThanIndex < lessThanIndex && greaterThanIndex != -1)) {
				if(greaterThanIndex + 1 != thanIndex) {
					throw new LOLError("THAN expected after BIGGR");
				}

				if(greaterThanIndex == 0 || greaterThanIndex == tokens.size() - 2) {
					throw new LOLError("Two arguments required for TIEMZ operation");
				}

				Expression expressionBefore = (tokens.get(greaterThanIndex - 1) instanceof Expression ? (Expression)tokens.get(greaterThanIndex - 1) : null);

				if(expressionBefore == null) {
					String expStringBefore = (String)tokens.get(greaterThanIndex - 1);

					if(expStringBefore.charAt(0) == '\"') {
						expressionBefore = new Value(new LOLString(expStringBefore.substring(1, expStringBefore.length() - 1)));
					} else {
						try {
							Double.parseDouble(expStringBefore);
						} catch(NumberFormatException e) {
							try {
								Long.parseLong(expStringBefore);
							} catch(NumberFormatException e2) {
								try {
									//expStringBefore = expStringBefore.toUpperCase();
									if(!expStringBefore.startsWith("0X")) {
										throw new NumberFormatException();
									}
									Long.parseLong(expStringBefore.replaceFirst("0X", ""), 16);
								} catch(NumberFormatException e3) {
									if(!expStringBefore.equals("YEZ") && !expStringBefore.equals("NO")) {
										expressionBefore = new VariableAndNoArgFunction(expStringBefore);
									}
								}
							}
						}
					}

					if(expressionBefore == null) {
						expressionBefore = new Value(LOLValue.valueOf(expStringBefore));
					}
				}

				Expression expressionAfter = (tokens.get(greaterThanIndex + 2) instanceof Expression ? (Expression)tokens.get(greaterThanIndex + 2) : null);

				if(expressionAfter == null) {
					String expStringAfter = (String)tokens.get(greaterThanIndex + 2);

					if(expStringAfter.charAt(0) == '\"') {
						expressionAfter = new Value(new LOLString(expStringAfter.substring(1, expStringAfter.length() - 1)));
					} else {
						try {
							Double.parseDouble(expStringAfter);
						} catch(NumberFormatException e) {
							try {
								Long.parseLong(expStringAfter);
							} catch(NumberFormatException e2) {
								try {
									//expStringAfter = expStringAfter.toUpperCase();
									if(!expStringAfter.startsWith("0X")) {
										throw new NumberFormatException();
									}
									Long.parseLong(expStringAfter.replaceFirst("0X", ""), 16);
								} catch(NumberFormatException e3) {
									if(!expStringAfter.equals("YEZ") && !expStringAfter.equals("NO")) {
										expressionAfter = new VariableAndNoArgFunction(expStringAfter);
									}
								}
							}
						}
					}

					if(expressionAfter == null) {
						expressionAfter = new Value(LOLValue.valueOf(expStringAfter));
					}
				}

				tokens.add(greaterThanIndex - 1, new GreaterThan(expressionBefore, expressionAfter));
				tokens.remove(greaterThanIndex);
				tokens.remove(greaterThanIndex);
				tokens.remove(greaterThanIndex);
				tokens.remove(greaterThanIndex);
				continue;
			}
		}

		while(tokens.contains("SAEM")) {
			int equalsIndex = tokens.indexOf("SAEM");

			if(tokens.indexOf("AS_SAEM") != equalsIndex + 1) {
				throw new LOLError("AS expected after SAEM");
			}

			if(equalsIndex == 0 || equalsIndex == tokens.size() - 2) {
				throw new LOLError("Two arguments required for AN operation");
			}

			Expression expressionBefore = (tokens.get(equalsIndex - 1) instanceof Expression ? (Expression)tokens.get(equalsIndex - 1) : null);

			if(expressionBefore == null) {
				String expStringBefore = (String)tokens.get(equalsIndex - 1);

				if(expStringBefore.charAt(0) == '\"') {
					expressionBefore = new Value(new LOLString(expStringBefore.substring(1, expStringBefore.length() - 1)));
				} else {
					try {
						Double.parseDouble(expStringBefore);
					} catch(NumberFormatException e) {
						try {
							Long.parseLong(expStringBefore);
						} catch(NumberFormatException e2) {
							try {
								//expStringBefore = expStringBefore.toUpperCase();
								if(!expStringBefore.startsWith("0X")) {
									throw new NumberFormatException();
								}
								Long.parseLong(expStringBefore.replaceFirst("0X", ""), 16);
							} catch(NumberFormatException e3) {
								if(!expStringBefore.equals("YEZ") && !expStringBefore.equals("NO")) {
									expressionBefore = new VariableAndNoArgFunction(expStringBefore);
								}
							}
						}
					}
				}

				if(expressionBefore == null) {
					expressionBefore = new Value(LOLValue.valueOf(expStringBefore));
				}
			}

			Expression expressionAfter = (tokens.get(equalsIndex + 2) instanceof Expression ? (Expression)tokens.get(equalsIndex + 2) : null);

			if(expressionAfter == null) {
				String expStringAfter = (String)tokens.get(equalsIndex + 2);

				if(expStringAfter.charAt(0) == '\"') {
					expressionAfter = new Value(new LOLString(expStringAfter.substring(1, expStringAfter.length() - 1)));
				} else {
					try {
						Double.parseDouble(expStringAfter);
					} catch(NumberFormatException e) {
						try {
							Long.parseLong(expStringAfter);
						} catch(NumberFormatException e2) {
							try {
								//expStringAfter = expStringAfter.toUpperCase();
								if(!expStringAfter.startsWith("0X")) {
									throw new NumberFormatException();
								}
								Long.parseLong(expStringAfter.replaceFirst("0X", ""), 16);
							} catch(NumberFormatException e3) {
								if(!expStringAfter.equals("YEZ") && !expStringAfter.equals("NO")) {
									expressionAfter = new VariableAndNoArgFunction(expStringAfter);
								}
							}
						}
					}
				}

				if(expressionAfter == null) {
					expressionAfter = new Value(LOLValue.valueOf(expStringAfter));
				}
			}

			tokens.add(equalsIndex - 1, new EqualTo(expressionBefore, expressionAfter));
			tokens.remove(equalsIndex);
			tokens.remove(equalsIndex);
			tokens.remove(equalsIndex);
			tokens.remove(equalsIndex);
			continue;
		}

		while(tokens.contains("AN")) {
			int andIndex = tokens.indexOf("AN");

			if(andIndex == 0 || andIndex == tokens.size() - 1) {
				throw new LOLError("Two arguments required for AN operation");
			}

			Expression expressionBefore = (tokens.get(andIndex - 1) instanceof Expression ? (Expression)tokens.get(andIndex - 1) : null);

			if(expressionBefore == null) {
				String expStringBefore = (String)tokens.get(andIndex - 1);

				if(expStringBefore.charAt(0) == '\"') {
					expressionBefore = new Value(new LOLString(expStringBefore.substring(1, expStringBefore.length() - 1)));
				} else {
					try {
						Double.parseDouble(expStringBefore);
					} catch(NumberFormatException e) {
						try {
							Long.parseLong(expStringBefore);
						} catch(NumberFormatException e2) {
							try {
								//expStringBefore = expStringBefore.toUpperCase();
								if(!expStringBefore.startsWith("0X")) {
									throw new NumberFormatException();
								}
								Long.parseLong(expStringBefore.replaceFirst("0X", ""), 16);
							} catch(NumberFormatException e3) {
								if(!expStringBefore.equals("YEZ") && !expStringBefore.equals("NO")) {
									expressionBefore = new VariableAndNoArgFunction(expStringBefore);
								}
							}
						}
					}
				}

				if(expressionBefore == null) {
					expressionBefore = new Value(LOLValue.valueOf(expStringBefore));
				}
			}

			Expression expressionAfter = (tokens.get(andIndex + 1) instanceof Expression ? (Expression)tokens.get(andIndex + 1) : null);

			if(expressionAfter == null) {
				String expStringAfter = (String)tokens.get(andIndex + 1);

				if(expStringAfter.charAt(0) == '\"') {
					expressionAfter = new Value(new LOLString(expStringAfter.substring(1, expStringAfter.length() - 1)));
				} else {
					try {
						Double.parseDouble(expStringAfter);
					} catch(NumberFormatException e) {
						try {
							Long.parseLong(expStringAfter);
						} catch(NumberFormatException e2) {
							try {
								//expStringAfter = expStringAfter.toUpperCase();
								if(!expStringAfter.startsWith("0X")) {
									throw new NumberFormatException();
								}
								Long.parseLong(expStringAfter.replaceFirst("0X", ""), 16);
							} catch(NumberFormatException e3) {
								if(!expStringAfter.equals("YEZ") && !expStringAfter.equals("NO")) {
									expressionAfter = new VariableAndNoArgFunction(expStringAfter);
								}
							}
						}
					}
				}

				if(expressionAfter == null) {
					expressionAfter = new Value(LOLValue.valueOf(expStringAfter));
				}
			}

			tokens.add(andIndex - 1, new LogicalAnd(expressionBefore, expressionAfter));
			tokens.remove(andIndex);
			tokens.remove(andIndex);
			tokens.remove(andIndex);
			continue;
		}

		while(tokens.contains("OR")) {
			int orIndex = tokens.indexOf("OR");

			if(orIndex == 0 || orIndex == tokens.size() - 1) {
				throw new LOLError("Two arguments required for AN operation");
			}

			Expression expressionBefore = (tokens.get(orIndex - 1) instanceof Expression ? (Expression)tokens.get(orIndex - 1) : null);

			if(expressionBefore == null) {
				String expStringBefore = (String)tokens.get(orIndex - 1);

				if(expStringBefore.charAt(0) == '\"') {
					expressionBefore = new Value(new LOLString(expStringBefore.substring(1, expStringBefore.length() - 1)));
				} else {
					try {
						Double.parseDouble(expStringBefore);
					} catch(NumberFormatException e) {
						try {
							Long.parseLong(expStringBefore);
						} catch(NumberFormatException e2) {
							try {
								//expStringBefore = expStringBefore.toUpperCase();
								if(!expStringBefore.startsWith("0X")) {
									throw new NumberFormatException();
								}
								Long.parseLong(expStringBefore.replaceFirst("0X", ""), 16);
							} catch(NumberFormatException e3) {
								if(!expStringBefore.equals("YEZ") && !expStringBefore.equals("NO")) {
									expressionBefore = new VariableAndNoArgFunction(expStringBefore);
								}
							}
						}
					}
				}

				if(expressionBefore == null) {
					expressionBefore = new Value(LOLValue.valueOf(expStringBefore));
				}
			}

			Expression expressionAfter = (tokens.get(orIndex + 1) instanceof Expression ? (Expression)tokens.get(orIndex + 1) : null);

			if(expressionAfter == null) {
				String expStringAfter = (String)tokens.get(orIndex + 1);

				if(expStringAfter.charAt(0) == '\"') {
					expressionAfter = new Value(new LOLString(expStringAfter.substring(1, expStringAfter.length() - 1)));
				} else {
					try {
						Double.parseDouble(expStringAfter);
					} catch(NumberFormatException e) {
						try {
							Long.parseLong(expStringAfter);
						} catch(NumberFormatException e2) {
							try {
								//expStringAfter = expStringAfter.toUpperCase();
								if(!expStringAfter.startsWith("0X")) {
									throw new NumberFormatException();
								}
								Long.parseLong(expStringAfter.replaceFirst("0X", ""), 16);
							} catch(NumberFormatException e3) {
								if(!expStringAfter.equals("YEZ") && !expStringAfter.equals("NO")) {
									expressionAfter = new VariableAndNoArgFunction(expStringAfter);
								}
							}
						}
					}
				}

				if(expressionAfter == null) {
					expressionAfter = new Value(LOLValue.valueOf(expStringAfter));
				}
			}

			tokens.add(orIndex - 1, new LogicalOr(expressionBefore, expressionAfter));
			tokens.remove(orIndex);
			tokens.remove(orIndex);
			tokens.remove(orIndex);
			continue;
		}

		if(tokens.size() != 1) {
			throw new LOLError("Unexpected symbol while parsing statement");
		}

		if(tokens.get(0) instanceof String) {
			String expString = (String)tokens.get(0);
			Expression expression = null;

			if(expString.charAt(0) == '\"') {
				expression = new Value(new LOLString(expString.substring(1, expString.length() - 1)));
			} else {
				try {
					Double.parseDouble(expString);
				} catch(NumberFormatException e) {
					try {
						Long.parseLong(expString);
					} catch(NumberFormatException e2) {
						try {
							//expString = expString.toUpperCase();
							if(!expString.startsWith("0X")) {
								throw new NumberFormatException();
							}
							Long.parseLong(expString.replaceFirst("0X", ""), 16);
						} catch(NumberFormatException e3) {
							if(!expString.equals("YEZ") && !expString.equals("NO")) {
								expression = new VariableAndNoArgFunction(expString);
							}
						}
					}
				}
			}

			if(expression == null) {
				expression = new Value(LOLValue.valueOf(expString));
			}

			tokens.add(0, expression);
			tokens.remove(1);
		}

		return (Expression)tokens.get(0);
	}

}
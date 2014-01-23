package org.objectivelol.vm;

import java.util.ArrayList;
import java.util.HashMap;

import org.objectivelol.lang.LOLBoolean;
import org.objectivelol.lang.LOLClass;
import org.objectivelol.lang.LOLError;
import org.objectivelol.lang.LOLFunction;
import org.objectivelol.lang.LOLNumber;
import org.objectivelol.lang.LOLObject;
import org.objectivelol.lang.LOLSource;
import org.objectivelol.lang.LOLValue;

public interface Expression {

	public LOLValue interpret(LOLObject owner, LOLFunction context, HashMap<String, ValueStruct> localVariables) throws LOLError, Return;

	public static class Return extends Throwable implements Expression {

		private static final long serialVersionUID = -984816694826243189L;

		private Expression right;
		private LOLValue value;

		public Return(Expression right) {
			this.right = right;
			this.value = null;
		}

		@Override
		public LOLValue interpret(LOLObject owner, LOLFunction context, HashMap<String, ValueStruct> localVariables) throws LOLError, Return {
			if(right != null) {
				value = right.interpret(owner, context, localVariables);
			}

			throw this;
		}

		public LOLValue getValue() {
			return value;
		}
	}

}

class Value implements Expression {

	private LOLValue value;

	public Value(LOLValue value) {
		this.value = value;
	}

	@Override
	public LOLValue interpret(LOLObject owner, LOLFunction context, HashMap<String, ValueStruct> localVariables) throws LOLError {
		return value;
	}

}

class VariableAndNoArgFunction implements Expression {

	private String name;

	public VariableAndNoArgFunction(String name) {
		this.name = name;
	}

	@Override
	public LOLValue interpret(LOLObject owner, LOLFunction context, HashMap<String, ValueStruct> localVariables) throws LOLError {
		ValueStruct vs = localVariables.get(name);

		if(vs == null) {
			if(owner == null) {
				vs = RuntimeEnvironment.getRuntime().getSource(context.getParentSource()).getGlobalVariable(name);
			} else {
				vs = owner.getVariable(name, context);

				if(vs == null) {
					vs = RuntimeEnvironment.getRuntime().getSource(context.getParentSource()).getGlobalVariable(name);
				}
			}
		}

		if(vs == null) {
			LOLFunction lf = null;

			if(owner == null) {
				lf = RuntimeEnvironment.getRuntime().getSource(context.getParentSource()).getGlobalFunction(name);
			} else {
				lf = owner.getFunction(name, context);

				if(lf == null) {
					lf = RuntimeEnvironment.getRuntime().getSource(context.getParentSource()).getGlobalFunction(name);
				}
			}

			if(lf == null) {
				throw new LOLError("The specified variable or function does not exist");
			}

			if(lf.isShared() == null || lf.isShared()) {
				return lf.execute(null, (LOLValue[])null);
			} else {
				return lf.execute(owner, (LOLValue[])null);
			}
		}

		return vs.getValue();
	}

}

class ArgFunction implements Expression {

	private String name;
	private ArrayList<Expression> arguments;

	public ArgFunction(String name, ArrayList<Expression> arguments) {
		this.name = name;
		this.arguments = arguments;
	}

	@Override
	public LOLValue interpret(LOLObject owner, LOLFunction context, HashMap<String, ValueStruct> localVariables) throws LOLError, Return {
		ArrayList<LOLValue> args = new ArrayList<LOLValue>(arguments.size());

		for(Expression e : arguments) {
			args.add(e.interpret(owner, context, localVariables));
		}

		LOLFunction lf = null;

		if(owner == null) {
			lf = RuntimeEnvironment.getRuntime().getSource(context.getParentSource()).getGlobalFunction(name);
		} else {
			lf = owner.getFunction(name, context);

			if(lf == null) {
				lf = RuntimeEnvironment.getRuntime().getSource(context.getParentSource()).getGlobalFunction(name);
			}
		}

		if(lf == null) {
			throw new LOLError("The specified function does not exist");
		}

		if(lf.isShared() == null || lf.isShared()) {
			return lf.execute(null, args.toArray(new LOLValue[args.size()]));
		} else {
			return lf.execute(owner, args.toArray(new LOLValue[args.size()]));
		}
	}

}

class Add implements Expression {

	private Expression left;
	private Expression right;

	public Add(Expression left, Expression right) {
		this.left = left;
		this.right = right;
	}

	@Override
	public LOLValue interpret(LOLObject owner, LOLFunction context, HashMap<String, ValueStruct> localVariables) throws LOLError, Return {
		return ((LOLNumber)left.interpret(owner, context, localVariables).cast(LOLNumber.TYPE_NAME)).add((LOLNumber)right.interpret(owner, context, localVariables).cast(LOLNumber.TYPE_NAME));
	}

}

class Subtract implements Expression {

	private Expression left;
	private Expression right;

	public Subtract(Expression left, Expression right) {
		this.left = left;
		this.right = right;
	}

	@Override
	public LOLValue interpret(LOLObject owner, LOLFunction context, HashMap<String, ValueStruct> localVariables) throws LOLError, Return {
		return ((LOLNumber)left.interpret(owner, context, localVariables).cast(LOLNumber.TYPE_NAME)).subtract((LOLNumber)right.interpret(owner, context, localVariables).cast(LOLNumber.TYPE_NAME));
	}

}

class Multiply implements Expression {

	private Expression left;
	private Expression right;

	public Multiply(Expression left, Expression right) {
		this.left = left;
		this.right = right;
	}

	@Override
	public LOLValue interpret(LOLObject owner, LOLFunction context, HashMap<String, ValueStruct> localVariables) throws LOLError, Return {
		return ((LOLNumber)left.interpret(owner, context, localVariables).cast(LOLNumber.TYPE_NAME)).multiply((LOLNumber)right.interpret(owner, context, localVariables).cast(LOLNumber.TYPE_NAME));
	}

}

class Divide implements Expression {

	private Expression left;
	private Expression right;

	public Divide(Expression left, Expression right) {
		this.left = left;
		this.right = right;
	}

	@Override
	public LOLValue interpret(LOLObject owner, LOLFunction context, HashMap<String, ValueStruct> localVariables) throws LOLError, Return {
		return ((LOLNumber)left.interpret(owner, context, localVariables).cast(LOLNumber.TYPE_NAME)).divide((LOLNumber)right.interpret(owner, context, localVariables).cast(LOLNumber.TYPE_NAME));
	}

}

class LogicalAnd implements Expression {

	private Expression left;
	private Expression right;
	
	public LogicalAnd(Expression left, Expression right) {
		this.left = left;
		this.right = right;
	}
	
	@Override
	public LOLValue interpret(LOLObject owner, LOLFunction context, HashMap<String, ValueStruct> localVariables) throws LOLError, Return {
		return LOLValue.valueOf(((LOLBoolean)left.interpret(owner, context, localVariables).cast(LOLBoolean.TYPE_NAME)).booleanValue() && ((LOLBoolean)right.interpret(owner, context, localVariables).cast(LOLBoolean.TYPE_NAME)).booleanValue());
	}
	
}

class LogicalOr implements Expression {

	private Expression left;
	private Expression right;
	
	public LogicalOr(Expression left, Expression right) {
		this.left = left;
		this.right = right;
	}
	
	@Override
	public LOLValue interpret(LOLObject owner, LOLFunction context, HashMap<String, ValueStruct> localVariables) throws LOLError, Return {
		return LOLValue.valueOf(((LOLBoolean)left.interpret(owner, context, localVariables).cast(LOLBoolean.TYPE_NAME)).booleanValue() || ((LOLBoolean)right.interpret(owner, context, localVariables).cast(LOLBoolean.TYPE_NAME)).booleanValue());
	}
	
}

class GreaterThan implements Expression {

	private Expression left;
	private Expression right;

	public GreaterThan(Expression left, Expression right) {
		this.left = left;
		this.right = right;
	}

	@Override
	public LOLValue interpret(LOLObject owner, LOLFunction context, HashMap<String, ValueStruct> localVariables) throws LOLError, Return {
		return ((LOLNumber)left.interpret(owner, context, localVariables).cast(LOLNumber.TYPE_NAME)).greaterThan((LOLNumber)right.interpret(owner, context, localVariables).cast(LOLNumber.TYPE_NAME));
	}

}

class LessThan implements Expression {

	private Expression left;
	private Expression right;

	public LessThan(Expression left, Expression right) {
		this.left = left;
		this.right = right;
	}

	@Override
	public LOLValue interpret(LOLObject owner, LOLFunction context, HashMap<String, ValueStruct> localVariables) throws LOLError, Return {
		return ((LOLNumber)left.interpret(owner, context, localVariables).cast(LOLNumber.TYPE_NAME)).lessThan((LOLNumber)right.interpret(owner, context, localVariables).cast(LOLNumber.TYPE_NAME));
	}

}

class EqualTo implements Expression {

	private Expression left;
	private Expression right;

	public EqualTo(Expression left, Expression right) {
		this.left = left;
		this.right = right;
	}

	@Override
	public LOLValue interpret(LOLObject owner, LOLFunction context, HashMap<String, ValueStruct> localVariables) throws LOLError, Return {
		return left.interpret(owner, context, localVariables).equalTo(right.interpret(owner, context, localVariables));
	}

}

class Cast implements Expression {

	private Expression left;
	private String targetType;

	public Cast(Expression left, String targetType) {
		this.left = left;
		this.targetType = targetType;
	}

	@Override
	public LOLValue interpret(LOLObject owner, LOLFunction context, HashMap<String, ValueStruct> localVariables) throws LOLError, Return {
		return left.interpret(owner, context, localVariables).cast(targetType);
	}

}

class DeclareVariable implements Expression {

	private String name;
	private String type;
	private boolean isLocked;
	private Expression right;

	public DeclareVariable(String name, String type, boolean isLocked, Expression right) {
		this.name = name;
		this.type = type;
		this.isLocked = isLocked;
		this.right = right;
	}

	@Override
	public LOLValue interpret(LOLObject owner, LOLFunction context, HashMap<String, ValueStruct> localVariables) throws LOLError, Return {
		ValueStruct vs = new ValueStruct(type, right.interpret(owner, context, localVariables), isLocked);

		if(localVariables.put(name, vs) != null) {
			throw new LOLError("Local variable identifier exists");
		}

		return null;
	}

}

class StatementBlock implements Expression {

	private ArrayList<Expression> statements;

	public StatementBlock(ArrayList<Expression> statements) {
		this.statements = statements;
	}

	@Override
	public LOLValue interpret(LOLObject owner, LOLFunction context, HashMap<String, ValueStruct> localVariables) throws LOLError, Return {
		for(Expression e : statements) {
			e.interpret(owner, context, localVariables);
		}

		return null;
	}

}

class WhileStatement implements Expression {

	private Expression condition;
	private Expression statements;

	public WhileStatement(Expression condition, Expression statements) {
		this.condition = condition;
		this.statements = statements;
	}

	@Override
	public LOLValue interpret(LOLObject owner, LOLFunction context, HashMap<String, ValueStruct> localVariables) throws LOLError, Return {
		while(condition.interpret(owner, context, localVariables).cast(LOLBoolean.TYPE_NAME).equalTo(LOLBoolean.YEZ).booleanValue()) {
			statements.interpret(owner, context, localVariables);
		}

		return null;
	}

}

class IfStatement implements Expression {

	private Expression condition;
	private Expression statements;

	public IfStatement(Expression condition, Expression statements) {
		this.condition = condition;
		this.statements = statements;
	}

	@Override
	public LOLValue interpret(LOLObject owner, LOLFunction context, HashMap<String, ValueStruct> localVariables) throws LOLError, Return {
		if(condition.interpret(owner, context, localVariables).cast(LOLBoolean.TYPE_NAME).equalTo(LOLBoolean.YEZ).booleanValue()) {
			statements.interpret(owner, context, localVariables);
		}

		return null;
	}

}

class SimpleAssignment implements Expression {

	private String name;
	private Expression right;

	public SimpleAssignment(String name, Expression right) {
		this.name = name;
		this.right = right;
	}

	@Override
	public LOLValue interpret(LOLObject owner, LOLFunction context, HashMap<String, ValueStruct> localVariables) throws LOLError, Return {
		ValueStruct vs = localVariables.get(name);

		if(vs == null) {
			if(owner == null) {
				vs = RuntimeEnvironment.getRuntime().getSource(context.getParentSource()).getGlobalVariable(name);
			} else {
				vs = owner.getVariable(name, context);

				if(vs == null) {
					vs = RuntimeEnvironment.getRuntime().getSource(context.getParentSource()).getGlobalVariable(name);
				}
			}
		}

		if(vs == null) {
			throw new LOLError("Variable not found");
		}
		
		if(vs.getIsLocked()) {
			throw new LOLError("Cannot assign value to LOCKD variable");
		}

		vs.setValue(right.interpret(owner, context, localVariables));

		return vs.getValue();
	}

}

class ComplexAssignment implements Expression {
	
	private String objectName;
	private String memberName;
	private Expression right;
	
	public ComplexAssignment(String objectName, String memberName, Expression right) {
		this.objectName = objectName;
		this.memberName = memberName;
		this.right = right;
	}

	@Override
	public LOLValue interpret(LOLObject owner, LOLFunction context, HashMap<String, ValueStruct> localVariables) throws LOLError, Return {
		LOLObject obj = null;
		
		ValueStruct vs = localVariables.get(objectName);

		if(vs == null) {
			if(owner == null) {
				vs = RuntimeEnvironment.getRuntime().getSource(context.getParentSource()).getGlobalVariable(objectName);
			} else {
				vs = owner.getVariable(objectName, context);

				if(vs == null) {
					vs = RuntimeEnvironment.getRuntime().getSource(context.getParentSource()).getGlobalVariable(objectName);
				}
			}
		}

		if(vs != null) {
			obj = (LOLObject)vs.getValue();
		}

		if(obj != null) {
			vs = obj.getVariable(memberName, context);
		} else {
			LOLClass lc = RuntimeEnvironment.getRuntime().getSource(context.getParentSource()).getGlobalClass(objectName);

			if(lc == null) {
				LOLSource ls = RuntimeEnvironment.getRuntime().getSource(objectName);

				if(ls == null) {
					throw new LOLError("Specified object not found");
				}

				vs = ls.getGlobalVariable(memberName);
			} else {
				vs = lc.getSharedVariable(memberName, context);
			}
		}
		
		if(vs == null) {
			throw new LOLError("Specified variable not found");
		}
		
		if(vs.getIsLocked()) {
			throw new LOLError("Cannot assign value to LOCKD variable");
		}
		
		vs.setValue(right.interpret(owner, context, localVariables));

		return vs.getValue();
	}
	
}

class MemberVariableAndNoArgFunction implements Expression {

	private String objectName;
	private String memberName;

	public MemberVariableAndNoArgFunction(String objectName, String memberName) {
		this.objectName = objectName;
		this.memberName = memberName;
	}

	@Override
	public LOLValue interpret(LOLObject owner, LOLFunction context, HashMap<String, ValueStruct> localVariables) throws LOLError, Return {
		LOLObject obj = null;

		ValueStruct vs = localVariables.get(objectName);

		if(vs == null) {
			if(owner == null) {
				vs = RuntimeEnvironment.getRuntime().getSource(context.getParentSource()).getGlobalVariable(objectName);
			} else {
				vs = owner.getVariable(objectName, context);

				if(vs == null) {
					vs = RuntimeEnvironment.getRuntime().getSource(context.getParentSource()).getGlobalVariable(objectName);
				}
			}
		}

		if(vs == null) {
			LOLFunction lf = null;

			if(owner == null) {
				lf = RuntimeEnvironment.getRuntime().getSource(context.getParentSource()).getGlobalFunction(objectName);
			} else {
				lf = owner.getFunction(objectName, context);

				if(lf == null) {
					lf = RuntimeEnvironment.getRuntime().getSource(context.getParentSource()).getGlobalFunction(objectName);
				}
			}

			if(lf != null) {
				if(lf.isShared() == null || lf.isShared()) {
					obj = (LOLObject)lf.execute(null, (LOLValue[])null);
				} else {
					obj = (LOLObject)lf.execute(owner, (LOLValue[])null);
				}
			}
		} else {
			obj = (LOLObject)vs.getValue();
		}

		if(obj != null) {
			vs = obj.getVariable(memberName, context);

			if(vs == null) {
				LOLFunction lf = obj.getFunction(memberName, context);

				if(lf == null) {
					throw new LOLError("Variable does not contain the specified member");
				}

				if(lf.isShared()) {
					return lf.execute(null, (LOLValue[])null);
				} else {
					return lf.execute(obj, (LOLValue[])null);
				}
			}

			return vs.getValue();
		} else {
			LOLClass lc = RuntimeEnvironment.getRuntime().getSource(context.getParentSource()).getGlobalClass(objectName);

			if(lc == null) {
				LOLSource ls = RuntimeEnvironment.getRuntime().getSource(objectName);

				if(ls == null) {
					throw new LOLError("Specified object not found");
				}

				vs = ls.getGlobalVariable(memberName);

				if(vs == null) {
					LOLFunction lf = ls.getGlobalFunction(memberName);

					if(lf == null) {
						throw new LOLError("Specified member not found");
					}

					return lf.execute(null, (LOLValue[])null);
				}

				return vs.getValue();
			} else {
				vs = lc.getSharedVariable(memberName, context);

				if(vs == null) {
					LOLFunction lf = lc.getSharedFunction(memberName, context);

					if(lf == null) {
						throw new LOLError("Specified member not found");
					}

					return lf.execute(null, (LOLValue[])null);
				}

				return vs.getValue();
			}
		}
	}

}

class MemberArgFunction implements Expression {

	private String objectName;
	private String functionName;
	private ArrayList<Expression> arguments;

	public MemberArgFunction(String objectName, String functionName, ArrayList<Expression> arguments) {
		this.objectName = objectName;
		this.functionName = functionName;
		this.arguments = arguments;
	}

	@Override
	public LOLValue interpret(LOLObject owner, LOLFunction context, HashMap<String, ValueStruct> localVariables) throws LOLError, Return {
		ArrayList<LOLValue> args = new ArrayList<LOLValue>(arguments.size());

		for(Expression e : arguments) {
			args.add(e.interpret(owner, context, localVariables));
		}

		LOLObject obj = null;

		ValueStruct vs = localVariables.get(objectName);

		if(vs == null) {
			if(owner == null) {
				vs = RuntimeEnvironment.getRuntime().getSource(context.getParentSource()).getGlobalVariable(objectName);
			} else {
				vs = owner.getVariable(objectName, context);

				if(vs == null) {
					vs = RuntimeEnvironment.getRuntime().getSource(context.getParentSource()).getGlobalVariable(objectName);
				}
			}
		}

		if(vs == null) {
			LOLFunction lf = null;

			if(owner == null) {
				lf = RuntimeEnvironment.getRuntime().getSource(context.getParentSource()).getGlobalFunction(objectName);
			} else {
				lf = owner.getFunction(objectName, context);

				if(lf == null) {
					lf = RuntimeEnvironment.getRuntime().getSource(context.getParentSource()).getGlobalFunction(objectName);
				}
			}

			if(lf != null) {
				if(lf.isShared() == null || lf.isShared()) {
					obj = (LOLObject)lf.execute(null, args.toArray(new LOLValue[args.size()]));
				} else {
					obj = (LOLObject)lf.execute(owner, args.toArray(new LOLValue[args.size()]));
				}
			}
		} else {
			obj = (LOLObject)vs.getValue();
		}

		if(obj != null) {
			LOLFunction lf = obj.getFunction(functionName, context);

			if(lf == null) {
				throw new LOLError("Variable does not contain the specified member");
			}

			if(lf.isShared()) {
				return lf.execute(null, args.toArray(new LOLValue[args.size()]));
			} else {
				return lf.execute(obj, args.toArray(new LOLValue[args.size()]));
			}
		} else {
			LOLClass lc = RuntimeEnvironment.getRuntime().getSource(context.getParentSource()).getGlobalClass(objectName);

			if(lc == null) {
				LOLSource ls = RuntimeEnvironment.getRuntime().getSource(objectName);

				if(ls == null) {
					throw new LOLError("Specified object not found");
				}

				LOLFunction lf = ls.getGlobalFunction(functionName);

				if(lf == null) {
					throw new LOLError("Specified member not found");
				}

				return lf.execute(null, args.toArray(new LOLValue[args.size()]));
			} else {
				LOLFunction lf = lc.getSharedFunction(functionName, context);

				if(lf == null) {
					throw new LOLError("Specified member not found");
				}

				return lf.execute(null, args.toArray(new LOLValue[args.size()]));
			}
		}
	}

}

class NewObject implements Expression {
	
	private String sourceName;
	private String className;
	
	public NewObject(String sourceName, String className) {
		this.sourceName = sourceName;
		this.className = className;
	}

	@Override
	public LOLValue interpret(LOLObject owner, LOLFunction context, HashMap<String, ValueStruct> localVariables) throws LOLError, Return {
		LOLClass lc = null;
		
		if(sourceName == null) {
			lc = RuntimeEnvironment.getRuntime().getSource(context.getParentSource()).getGlobalClass(className);
		} else {
			lc = RuntimeEnvironment.getRuntime().getSource(sourceName).getGlobalClass(className);
		}
		
		if(lc == null) {
			throw new LOLError("Specified class not found");
		}
		
		return lc.constructInstance();
	}
	
}
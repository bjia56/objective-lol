package org.objectivelol.vm;

import org.objectivelol.lang.LOLBoolean;
import org.objectivelol.lang.LOLClass;
import org.objectivelol.lang.LOLError;
import org.objectivelol.lang.LOLFunction;
import org.objectivelol.lang.LOLObject;
import org.objectivelol.lang.LOLValue;

class LOLObjectRuntimeWrapper extends LOLObject {

	private String sourceName;
	private String className;
	private LOLObject obj;

	public LOLObjectRuntimeWrapper(String sourceName, String className) {
		super(null, null, null);

		this.sourceName = sourceName;
		this.className = className;
		this.obj = null;
	}

	@Override
	public LOLValue cast(String type) throws LOLError {
		if(obj == null) {
			instantiate();
		}

		return obj.cast(type);
	}

	@Override
	public String getTypeName() {
		return className;
	}

	@Override
	public LOLBoolean equalTo(LOLValue other) throws LOLError {
		if(obj == null) {
			instantiate();
		}

		return obj.equalTo(other);
	}

	public LOLFunction getFunction(String name, LOLFunction context) throws LOLError {
		if(obj == null) {
			instantiate();
		}

		return obj.getFunction(name, context);
	}

	public ValueStruct getVariable(String name, LOLFunction context) throws LOLError {
		if(obj == null) {
			instantiate();
		}

		return obj.getVariable(name, context);
	}

	private void instantiate() throws LOLError {
		LOLClass lc = null;

		lc = RuntimeEnvironment.getRuntime().getSource(sourceName).getGlobalClass(className);

		if(lc == null) {
			throw new LOLError("Specified class not found when attempting post-parsing instantiation");
		}

		obj = lc.constructInstance();
	}
	
	@Override
	public LOLValue copy() throws LOLError {
		if(obj == null) {
			return new LOLObjectRuntimeWrapper(sourceName, className);
		}
		
		return obj.copy();
	}

}
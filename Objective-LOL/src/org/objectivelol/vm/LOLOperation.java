package org.objectivelol.vm;

import org.objectivelol.lang.LOLFunction;
import org.objectivelol.lang.LOLValue;

public abstract class LOLOperation {

	public abstract LOLValue execute(LOLFunction context);
	
}

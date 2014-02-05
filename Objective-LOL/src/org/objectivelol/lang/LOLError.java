package org.objectivelol.lang;

/**
 * Class to represent an error of any type in
 * Objective-LOL. Does not map to any feature
 * of the language; instead, this is thrown by the
 * virtual machine when any parsing or runtime
 * errors arise.
 * 
 * @author Brett Jia
 */
public class LOLError extends Exception {

	public static final long serialVersionUID = 685447603884690089L;

	/**
	 * Constructor for the LOLError class.
	 * 
	 * @param s
	 * A String representing what error was present
	 * and caused this LOLError to be thrown.
	 */
	public LOLError(String s) {
		super(s);
	}
	
}

import org.objectivelol.lang.LOLError;
import org.objectivelol.lang.LOLValue;
import org.objectivelol.vm.RuntimeEnvironment;

public class MainClass {

	/**
	 * @param args
	 * @throws LOLError 
	 */
	public static void main(String[] args) throws LOLError {
		RuntimeEnvironment re = RuntimeEnvironment.getRuntime();

		re.loadSource("test.lol");

		LOLValue result = re.getSource("TEST").getGlobalFunction("FIBONACCI").execute(null, LOLValue.valueOf("20"));
		re.getSource("STDIO").getGlobalFunction("COMPLAIN").execute(null, result);
		//re.getSource("STDIO").getGlobalFunction("VISIBLE").execute(null, re.getSource("STDIO").getGlobalFunction("GIMMEH").execute(null));
		result = re.getSource("TEST").getGlobalFunction("MAIN").execute(null);
		re.execute();
	}

}

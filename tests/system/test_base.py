from nbpbeat import BaseTest

import os


class Test(BaseTest):

    def test_base(self):
        """
        Basic test with exiting Nbpbeat normally
        """
        self.render_config_template(
            path=os.path.abspath(self.working_dir) + "/log/*"
        )

        nbpbeat_proc = self.start_beat()
        self.wait_until(lambda: self.log_contains("nbpbeat is running"))
        exit_code = nbpbeat_proc.kill_and_wait()
        assert exit_code == 0

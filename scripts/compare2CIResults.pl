#!/usr/bin/perl  
use Pod::Usage;
use Getopt::Long;
#######################################################################################################################################################################################
#This script allows you to compare 2 CI regression results and generate a report for tests, which are in the 1st file; not in the 2nd file.
## @author : tlin
#Steps to use the script:
#1. download the script
#2. copy a CI regression result from Jenkins website and save it into a file ( see ciresults_1 as an example)
#3. copy another CI regression result from Jenkins website and save it into a file ( see ciresults_2 as an example)
#4. pls make sure that the data format of your CI resgression results match the examples, ciresults_1 and ciresults_2
#5 ./compare2CIResults.pl ciresults_1 ciresults_2 report
#6. after step6, the file, report, will have all tests, which are in 1st file; but not in the 2nd file.
######################################################################################################################################################################################################


$file1name = $ARGV[0];
$file2name = $ARGV[1]; 
$reportFile = $ARGV[2];

my $help = 0;
GetOptions('help|?' => \$help,);
pod2usage( { -message => "This script allows you to compare 2 CI regression results and generate a report for tests, \n which are in the 1st file; not in the 2nd file. \n\nUsage: the script takes 3 arguments. Here is the example command \n ./compare2CIResults.pl ciresults_1 ciresults_2 report \n" ,
-exitval => 1 } )if $help;

die "This script allows you to compare 2 CI regression results and generate a report for tests, \n which are in the 1st file; not in the 2nd file. \n\nUsage: the script takes 3 arguments. Here is the example command \n ./compare2CIResults.pl ciresults_1 ciresults_2 report \n" if @ARGV != 3;

%results1;
%results2; 
%results1WithCompletedTestCaseName;
%results2WithCompletedTestCaseName; 

open(INFO, $file1name) or die("Could not open  file.");

foreach $line (<INFO>)  {  
	@info = split ('\.', $line);
	@temp = split (" ", $line);
	#print $info[1] \n; 
	$feature = trim("com.$info[1].$info[2]");  
        if (trim($line) =~ /^>/) {
                $completeTestCaseName = trim($temp[1]);
        } else {
                $completeTestCaseName = trim($temp[0]);
        }
	$results1WithCompletedTestCaseName {$completeTestCaseName} = 1;

	if (! exists $results1{$feature}) {
		$results1{$feature} = 1;
	} else {
		$results1{$feature} = $results1{$feature} +1;
	}

}
close(INFO);

open(INFO, $file2name) or die("Could not open  file.");
foreach $line (<INFO>)  {  
	@info = split ('\.', $line);
	@temp = split (" ", $line);
	#print $info[1] \n; 
	$feature = trim("com.$info[1].$info[2]");  
        if (trim($line) =~ /^>/) {
                $completeTestCaseName = trim($temp[1]);
        } else {
                $completeTestCaseName = trim($temp[0]);
        }
	$results2WithCompletedTestCaseName {$completeTestCaseName} = 1;
	if (! exists $results2{$feature}) {
		$results2{$feature} = 1;
	} else {
		$results2{$feature} = $results2{$feature} +1;
	}

}
close(INFO);

system("echo \"Here is the list of tests, which are in the 1st file; but not the 2nd file \n \" > $reportFile  ");
foreach $testcaseName (keys %results1WithCompletedTestCaseName) {
	if (! exists $results2WithCompletedTestCaseName{$testcaseName}) {
		system ("echo \"$testcaseName\" >> $reportFile ");
	}
} 
 
sub trim($)
{
	my $string = shift;
	$string =~ s/^\s+//;
	$string =~ s/\s+$//;
	return $string;
}


sub printHashTable { 
	my $hashRef = shift;
	my %subHash = %{$hashRef};
	foreach $key (keys %subHash) {
		$value = $subHash{$key};
		print "$key ==> $value \n";
	}
}






!define PRODUCT_NAME "CloudOps Agent"
!define PRODUCT_NAME_OTHER "CloudOps"
!define PRODUCT_PUBLISHER "CloudOps"
!define PRODUCT_WEB_SITE "https://cloudops.thinkcloud.lenovo.com/"
!define PRODUCT_BIN_REGKEY "Software\Microsoft\Windows\CurrentVersion\App Paths\agent.exe"
!define PRODUCT_UNINST_KEY "Software\Microsoft\Windows\CurrentVersion\Uninstall\${PRODUCT_NAME}"
!define PRODUCT_UNINST_KEY_OTHER "Software\Microsoft\Windows\CurrentVersion\Uninstall\${PRODUCT_NAME_OTHER}"
!define PRODUCT_UNINST_ROOT_KEY "HKLM"
!define OUTFILE "CloudOps-Agent-${CPUARCH}-Setup.exe"

# Import Libraries
!include "MUI2.nsh"
!include "nsDialogs.nsh"
!include "LogicLib.nsh"
!include "FileFunc.nsh"
!include "StrFunc.nsh"
!include "x64.nsh"
!include "WinMessages.nsh"
!include "WinVer.nsh"
${StrLoc}

!ifdef Version
    !define PRODUCT_VERSION "${Version}"
!else
    !define PRODUCT_VERSION "Undefined Version"
!endif

!if "$%PROCESSOR_ARCHITECTURE%" == "AMD64"
    !define CPUARCH "AMD64"
!else if "$%PROCESSOR_ARCHITEW6432%" == "AMD64"
    !define CPUARCH "AMD64"
!else
    !define CPUARCH "x86"
!endif

# Part of the Trim function for Strings
!define Trim "!insertmacro Trim"
!macro Trim ResultVar String
    Push "${String}"
    Call Trim
    Pop "${ResultVar}"
!macroend

# Part of the Explode function for Strings
!define Explode "!insertmacro Explode"
!macro Explode Length Separator String
    Push    `${Separator}`
    Push    `${String}`
    Call    Explode
    Pop     `${Length}`
!macroend


###############################################################################
# Configure Pages, Ordering, and Configuration
###############################################################################
!define MUI_ABORTWARNING
!define MUI_ICON "salt.ico"
!define MUI_UNICON "salt.ico"
!define MUI_WELCOMEFINISHPAGE_BITMAP "panel.bmp"
!define MUI_UNWELCOMEFINISHPAGE_BITMAP "panel.bmp"


# Welcome page
!insertmacro MUI_PAGE_WELCOME

# License page
!insertmacro MUI_PAGE_LICENSE "LICENSE.txt"

# Configure Minion page
Page custom pageMinionConfig pageMinionConfig_Leave


# Instfiles page
!insertmacro MUI_PAGE_INSTFILES

# Finish page (Customized)
!define MUI_PAGE_CUSTOMFUNCTION_SHOW pageFinish_Show
!define MUI_PAGE_CUSTOMFUNCTION_LEAVE pageFinish_Leave
!insertmacro MUI_PAGE_FINISH

# Uninstaller pages
!insertmacro MUI_UNPAGE_INSTFILES

# Language files
!insertmacro MUI_LANGUAGE "English"


###############################################################################
# Custom Dialog Box Variables
###############################################################################
Var Dialog
Var Label
Var CheckBox_Minion_Start
Var CheckBox_Minion_Start_Delayed
Var ApiKey
Var ApiKey_State
Var Endpoint
Var Endpoint_State
Var StartMinion
Var StartMinionDelayed
Var DeleteInstallDir
Var ConfigWriteMinion
Var ConfigWriteMaster


###############################################################################
# Minion Settings Dialog Box
###############################################################################
Function pageMinionConfig

    # Set Page Title and Description
    !insertmacro MUI_HEADER_TEXT "Settings" "Set the Agent ApiKey and Endpoint"
    nsDialogs::Create 1018
    Pop $Dialog

    ${If} $Dialog == error
        Abort
    ${EndIf}

    # ApiKey Dialog Control
    ${NSD_CreateLabel} 0 0 100% 12u "ApiKey:"
    Pop $Label

    ${NSD_CreateText} 0 13u 100% 12u $ApiKey_State
    Pop $ApiKey

    # Minion ID Dialog Control
    ${NSD_CreateLabel} 0 30u 100% 12u "Endpoint:"
    Pop $Label

    ${NSD_CreateText} 0 43u 100% 12u $Endpoint_State
    Pop $Endpoint

    nsDialogs::Show

FunctionEnd



# File Picker Definitions
!define OFN_FILEMUSTEXIST 0x00001000
!define OFN_DONTADDTOREC 0x02000000
!define OPENFILENAME_SIZE_VERSION_400 76
!define OPENFILENAME 'i,i,i,i,i,i,i,i,i,i,i,i,i,i,&i2,&i2,i,i,i,i'

Function pageMinionConfig_Leave

    # Save the State
    ${NSD_GetText} $ApiKey $ApiKey_State
    ${NSD_GetText} $Endpoint $Endpoint_State

FunctionEnd


###############################################################################
# Custom Finish Page
###############################################################################
Function pageFinish_Show

    # Imports so the checkboxes will show up
    !define SWP_NOSIZE 0x0001
    !define SWP_NOMOVE 0x0002
    !define HWND_TOP 0x0000

    # Create Start Minion Checkbox
    ${NSD_CreateCheckbox} 120u 90u 100% 12u "&Start cloudops-agent"
    Pop $CheckBox_Minion_Start
    SetCtlColors $CheckBox_Minion_Start "" "ffffff"
    # This command required to bring the checkbox to the front
    System::Call "User32::SetWindowPos(i, i, i, i, i, i, i) b ($CheckBox_Minion_Start, ${HWND_TOP}, 0, 0, 0, 0, ${SWP_NOSIZE}|${SWP_NOMOVE})"

    # Create Start Minion Delayed ComboBox
    ${NSD_CreateCheckbox} 130u 102u 100% 12u "&Delayed Start"
    Pop $CheckBox_Minion_Start_Delayed
    SetCtlColors $CheckBox_Minion_Start_Delayed "" "ffffff"
    # This command required to bring the checkbox to the front
    System::Call "User32::SetWindowPos(i, i, i, i, i, i, i) b ($CheckBox_Minion_Start_Delayed, ${HWND_TOP}, 0, 0, 0, 0, ${SWP_NOSIZE}|${SWP_NOMOVE})"

    # Load current settings for Minion
    ${If} $StartMinion == 1
        ${NSD_Check} $CheckBox_Minion_Start
    ${EndIf}

    # Load current settings for Minion Delayed
    ${If} $StartMinionDelayed == 1
        ${NSD_Check} $CheckBox_Minion_Start_Delayed
    ${EndIf}

FunctionEnd


Function pageFinish_Leave

    # Assign the current checkbox states
    ${NSD_GetState} $CheckBox_Minion_Start $StartMinion
    ${NSD_GetState} $CheckBox_Minion_Start_Delayed $StartMinionDelayed

FunctionEnd


###############################################################################
# Installation Settings
###############################################################################
Name "${PRODUCT_NAME} ${PRODUCT_VERSION}"
OutFile "${OutFile}"
InstallDir "c:\cloudops-agent"
InstallDirRegKey HKLM "${PRODUCT_DIR_REGKEY}" ""
ShowInstDetails show
ShowUnInstDetails show

Section -copy_prereqs
    # Copy prereqs to the Plugins Directory
    # These files will be vcredist 2008 and KB2999226 for Win8.1 and below
    # These files are downloaded by build_pkg.bat
    # This directory gets removed upon completion
    SetOutPath "$PLUGINSDIR\"
    File /r "..\prereqs\"
SectionEnd

# Check and install the Windows 10 Universal C Runtime (KB2999226)
# ucrt is needed on Windows 8.1 and lower
# They are installed as a Microsoft Update package (.msu)
# ucrt for Windows 8.1 RT is only available via Windows Update
Section -install_ucrt

    Var /GLOBAL MsuPrefix
    Var /GLOBAL MsuFileName

    # Get the Major.Minor version Number
    # Windows 10 introduced CurrentMajorVersionNumber
    ReadRegStr $R0 HKLM "SOFTWARE\Microsoft\Windows NT\CurrentVersion" \
        CurrentMajorVersionNumber

    # Windows 10/2016 will return a value here, skip to the end if returned
    StrCmp $R0 '' lbl_needs_ucrt 0

    # Found Windows 10
    detailPrint "KB2999226 does not apply to this machine"
    goto lbl_done

    lbl_needs_ucrt:
    # UCRT only needed on Windows Server 2012R2/Windows 8.1 and below
    # The first ReadRegStr command above should have skipped to lbl_done if on
    # Windows 10 box

    # Is the update already installed
    ClearErrors

    # Use WMI to check if it's installed
    detailPrint "Checking for existing KB2999226 installation"
    nsExec::ExecToStack 'cmd /q /c wmic qfe get hotfixid | findstr "^KB2999226"'
    # Clean up the stack
    Pop $R0 # Gets the ErrorCode
    Pop $R1 # Gets the stdout, which should be KB2999226 if it's installed

    # If it returned KB2999226 it's already installed
    StrCmp $R1 'KB2999226' lbl_done

    detailPrint "KB2999226 not found"

    # All lower versions of Windows
    ReadRegStr $R0 HKLM "SOFTWARE\Microsoft\Windows NT\CurrentVersion" \
        CurrentVersion

    # Get the name of the .msu file based on the value of $R0
    ${Switch} "$R0"
        ${Case} "6.3"
            StrCpy $MsuPrefix "Windows8.1"
            ${break}
        ${Case} "6.2"
            StrCpy $MsuPrefix "Windows8-RT"
            ${break}
        ${Case} "6.1"
            StrCpy $MsuPrefix "Windows6.1"
            ${break}
        ${Case} "6.0"
            StrCpy $MsuPrefix "Windows6.0"
            ${break}
    ${EndSwitch}

    # Use RunningX64 here to get the Architecture for the system running the installer
    # CPUARCH is defined when the installer is built and is based on the machine that
    # built the installer, not the target system as we need here.
    ${If} ${RunningX64}
        StrCpy $MsuFileName "$MsuPrefix-KB2999226-x64.msu"
    ${Else}
        StrCpy $MsuFileName "$MsuPrefix-KB2999226-x86.msu"
    ${EndIf}

    ClearErrors

    detailPrint "Installing KB2999226 using file $MsuFileName"
    nsExec::ExecToStack 'cmd /c wusa "$PLUGINSDIR\$MsuFileName" /quiet /norestart'
    # Clean up the stack
    Pop $R0  # Get Error
    Pop $R1  # Get stdout
    ${IfNot} $R0 == 0
        detailPrint "error: $R0"
        detailPrint "output: $R2"
        Sleep 3000
    ${Else}
        detailPrint "KB2999226 installed successfully"
    ${EndIf}

    lbl_done:

SectionEnd


Section "MainSection" SEC01

    SetOutPath "$INSTDIR\"
    SetOverwrite off
    File /r "..\buildenv\"

SectionEnd


Function .onInit

    Call parseCommandLineSwitches

    # If custom config passed, verify its existence before continuing so we
    # don't uninstall an existing installation and then fail
    # Check for existing installation
    ReadRegStr $R0 HKLM \
        "Software\Microsoft\Windows\CurrentVersion\Uninstall\${PRODUCT_NAME}" \
        "UninstallString"
    StrCmp $R0 "" checkOther
    # Found existing installation, prompt to uninstall
    MessageBox MB_OKCANCEL|MB_ICONEXCLAMATION \
        "${PRODUCT_NAME} is already installed.$\n$\n\
        Click `OK` to remove the existing installation." \
        /SD IDOK IDOK uninst
    Abort

    checkOther:
        # Check for existing installation of full salt
        ReadRegStr $R0 HKLM \
            "Software\Microsoft\Windows\CurrentVersion\Uninstall\${PRODUCT_NAME_OTHER}" \
            "UninstallString"
        StrCmp $R0 "" skipUninstall
        # Found existing installation, prompt to uninstall
        MessageBox MB_OKCANCEL|MB_ICONEXCLAMATION \
            "${PRODUCT_NAME_OTHER} is already installed.$\n$\n\
            Click `OK` to remove the existing installation." \
            /SD IDOK IDOK uninst
        Abort

    uninst:

        # Get current Silent status
        StrCpy $R0 0
        ${If} ${Silent}
            StrCpy $R0 1
        ${EndIf}

        # Turn on Silent mode
        SetSilent silent

        # Don't remove all directories
        StrCpy $DeleteInstallDir 0

        # Uninstall silently
        Call uninstallSalt

        # Set it back to Normal mode, if that's what it was before
        ${If} $R0 == 0
            SetSilent normal
        ${EndIf}

    skipUninstall:


FunctionEnd


Section -Post

    WriteUninstaller "$INSTDIR\uninst.exe"

    # Uninstall Registry Entries
    WriteRegStr ${PRODUCT_UNINST_ROOT_KEY} "${PRODUCT_UNINST_KEY}" \
        "DisplayName" "$(^Name)"
    WriteRegStr ${PRODUCT_UNINST_ROOT_KEY} "${PRODUCT_UNINST_KEY}" \
        "UninstallString" "$INSTDIR\uninst.exe"
    WriteRegStr ${PRODUCT_UNINST_ROOT_KEY} "${PRODUCT_UNINST_KEY}" \
        "DisplayIcon" "$INSTDIR\salt.ico"
    WriteRegStr ${PRODUCT_UNINST_ROOT_KEY} "${PRODUCT_UNINST_KEY}" \
        "DisplayVersion" "${PRODUCT_VERSION}"
    WriteRegStr ${PRODUCT_UNINST_ROOT_KEY} "${PRODUCT_UNINST_KEY}" \
        "URLInfoAbout" "${PRODUCT_WEB_SITE}"
    WriteRegStr ${PRODUCT_UNINST_ROOT_KEY} "${PRODUCT_UNINST_KEY}" \
        "Publisher" "${PRODUCT_PUBLISHER}"
    WriteRegStr HKLM "SYSTEM\CurrentControlSet\services\cloudops-agent" \
        "DependOnService" "nsi"

    # Set the estimated size
    ${GetSize} "$INSTDIR" "/S=OK" $0 $1 $2
    IntFmt $0 "0x%08X" $0
    WriteRegDWORD ${PRODUCT_UNINST_ROOT_KEY} "${PRODUCT_UNINST_KEY}" \
        "EstimatedSize" "$0"

    # Commandline Registry Entries
    WriteRegStr HKLM "${PRODUCT_BIN_REGKEY}" "Path" "$INSTDIR\"

    # Register the Salt-Minion Service
    nsExec::Exec "$INSTDIR\ssm.exe install cloudops-agent $INSTDIR\agent.exe"
    nsExec::Exec "$INSTDIR\ssm.exe set cloudops-agent AppParameters run"
    nsExec::Exec "$INSTDIR\ssm.exe set cloudops-agent AppDirectory $INSTDIR"
    nsExec::Exec "$INSTDIR\ssm.exe set cloudops-agent Description CloudOps Agent"
    nsExec::Exec "$INSTDIR\ssm.exe set cloudops-agent Start SERVICE_AUTO_START"
    nsExec::Exec "$INSTDIR\ssm.exe set cloudops-agent AppStopMethodConsole 24000"
    nsExec::Exec "$INSTDIR\ssm.exe set cloudops-agent AppStopMethodWindow 2000"

    Call updateMinionConfig

    Push "C:\cloudops-agent"
    Call AddToPath

SectionEnd


Function .onInstSuccess

    # If StartMinionDelayed is 1, then set the service to start delayed
    ${If} $StartMinionDelayed == 1
        nsExec::Exec "$INSTDIR\ssm.exe set cloudops-agent Start SERVICE_DELAYED_AUTO_START"
    ${EndIf}

    # If start-minion is 1, then start the service
    ${If} $StartMinion == 1
        nsExec::Exec 'net start cloudops-agent'
    ${EndIf}

FunctionEnd


Function un.onInit

    # Load the parameters
    ${GetParameters} $R0

    # Uninstaller: Remove Installation Directory
    ClearErrors
    ${GetOptions} $R0 "/delete-install-dir" $R1
    IfErrors delete_install_dir_not_found
        StrCpy $DeleteInstallDir 1
    delete_install_dir_not_found:

    MessageBox MB_ICONQUESTION|MB_YESNO|MB_DEFBUTTON2 \
        "Are you sure you want to completely remove $(^Name) and all of its components?" \
        /SD IDYES IDYES +2
    Abort

FunctionEnd


Section Uninstall

    Call un.uninstallSalt

    # Remove C:\cloudops-agent from the Path
    Push "C:\cloudops-agent"
    Call un.RemoveFromPath

SectionEnd


!macro uninstallSalt un
Function ${un}uninstallSalt

    # Make sure we're in the right directory
    StrCpy $INSTDIR "C:\cloudops-agent"

    # Stop and Remove cloudops-qgent service
    nsExec::Exec 'net stop cloudops-agent'
    nsExec::Exec 'sc delete cloudops-agent'

    # Remove files
    Delete "$INSTDIR\uninst.exe"
    Delete "$INSTDIR\ssm.exe"
    RMDir /r "$INSTDIR"

    # Remove Registry entries
    DeleteRegKey ${PRODUCT_UNINST_ROOT_KEY} "${PRODUCT_BIN_REGKEY}"

    # Automatically close when finished
    SetAutoClose true

    # Prompt to remove the Installation directory
    ${IfNot} $DeleteInstallDir == 1
        MessageBox MB_ICONQUESTION|MB_YESNO|MB_DEFBUTTON2 \
            "Would you like to completely remove $INSTDIR and all of its contents?" \
            /SD IDNO IDNO finished
    ${EndIf}

    # Make sure you're not removing Program Files
    ${If} $INSTDIR != 'Program Files'
    ${AndIf} $INSTDIR != 'Program Files (x86)'
        RMDir /r "$INSTDIR"
    ${EndIf}

    finished:

FunctionEnd
!macroend


!insertmacro uninstallSalt ""
!insertmacro uninstallSalt "un."


Function un.onUninstSuccess
    HideWindow
    MessageBox MB_ICONINFORMATION|MB_OK \
        "$(^Name) was successfully removed from your computer." \
        /SD IDOK
FunctionEnd

#------------------------------------------------------------------------------
# Trim Function
# - Trim whitespace from the beginning and end of a string
# - Trims spaces, \r, \n, \t
#
# Usage:
#   Push " some string "  ; String to Trim
#   Call Trim
#   Pop $0                ; Trimmed String: "some string"
#
#   or
#
#   ${Trim} $0 $1   ; Trimmed String, String to Trim
#------------------------------------------------------------------------------
Function Trim

    Exch $R1 # Original string
    Push $R2

    Loop:
        StrCpy $R2 "$R1" 1
        StrCmp "$R2" " " TrimLeft
        StrCmp "$R2" "$\r" TrimLeft
        StrCmp "$R2" "$\n" TrimLeft
        StrCmp "$R2" "$\t" TrimLeft
        GoTo Loop2
    TrimLeft:
        StrCpy $R1 "$R1" "" 1
        Goto Loop

    Loop2:
        StrCpy $R2 "$R1" 1 -1
        StrCmp "$R2" " " TrimRight
        StrCmp "$R2" "$\r" TrimRight
        StrCmp "$R2" "$\n" TrimRight
        StrCmp "$R2" "$\t" TrimRight
        GoTo Done
    TrimRight:
        StrCpy $R1 "$R1" -1
        Goto Loop2

    Done:
        Pop $R2
        Exch $R1

FunctionEnd


#------------------------------------------------------------------------------
# Explode Function
# - Splits a string based off the passed separator
# - Each item in the string is pushed to the stack
# - The last item pushed to the stack is the length of the array
#
# Usage:
#   Push ","                    ; Separator
#   Push "string,to,separate"   ; String to explode
#   Call Explode
#   Pop $0                      ; Number of items in the array
#
#   or
#
#   ${Explode} $0 $1 $2         ; Length, Separator, String
#------------------------------------------------------------------------------
Function Explode
    # Initialize variables
    Var /GLOBAL explString
    Var /GLOBAL explSeparator
    Var /GLOBAL explStrLen
    Var /GLOBAL explSepLen
    Var /GLOBAL explOffset
    Var /GLOBAL explTmp
    Var /GLOBAL explTmp2
    Var /GLOBAL explTmp3
    Var /GLOBAL explArrCount

    # Get input from user
    Pop $explString
    Pop $explSeparator

    # Calculates initial values
    StrLen $explStrLen $explString
    StrLen $explSepLen $explSeparator
    StrCpy $explArrCount 1

    ${If} $explStrLen <= 1             #   If we got a single character
    ${OrIf} $explSepLen > $explStrLen  #   or separator is larger than the string,
        Push    $explString            #   then we return initial string with no change
        Push    1                      #   and set array's length to 1
        Return
    ${EndIf}

    # Set offset to the last symbol of the string
    StrCpy $explOffset $explStrLen
    IntOp  $explOffset $explOffset - 1

    # Clear temp string to exclude the possibility of appearance of occasional data
    StrCpy $explTmp   ""
    StrCpy $explTmp2  ""
    StrCpy $explTmp3  ""

    # Loop until the offset becomes negative
    ${Do}
        # If offset becomes negative, it is time to leave the function
        ${IfThen} $explOffset == -1 ${|} ${ExitDo} ${|}

        # Remove everything before and after the searched part ("TempStr")
        StrCpy $explTmp $explString $explSepLen $explOffset

        ${If} $explTmp == $explSeparator
            # Calculating offset to start copy from
            IntOp   $explTmp2 $explOffset + $explSepLen    # Offset equals to the current offset plus length of separator
            StrCpy  $explTmp3 $explString "" $explTmp2

            Push    $explTmp3                              # Throwing array item to the stack
            IntOp   $explArrCount $explArrCount + 1        # Increasing array's counter

            StrCpy  $explString $explString $explOffset 0  # Cutting all characters beginning with the separator entry
            StrLen  $explStrLen $explString
        ${EndIf}

        ${If} $explOffset = 0           # If the beginning of the line met and there is no separator,
                                        # copying the rest of the string
            ${If} $explSeparator == ""  # Fix for the empty separator
                IntOp   $explArrCount   $explArrCount - 1
            ${Else}
                Push    $explString
            ${EndIf}
        ${EndIf}

        IntOp   $explOffset $explOffset - 1
    ${Loop}

    Push $explArrCount
FunctionEnd


#------------------------------------------------------------------------------
# StrStr Function
# - find substring in a string
#
# Usage:
#   Push "this is some string"
#   Push "some"
#   Call StrStr
#   Pop $0 # "some string"
#------------------------------------------------------------------------------
!macro StrStr un
Function ${un}StrStr

    Exch $R1 # $R1=substring, stack=[old$R1,string,...]
    Exch     #                stack=[string,old$R1,...]
    Exch $R2 # $R2=string,    stack=[old$R2,old$R1,...]
    Push $R3 # $R3=strlen(substring)
    Push $R4 # $R4=count
    Push $R5 # $R5=tmp
    StrLen $R3 $R1 # Get the length of the Search String
    StrCpy $R4 0 # Set the counter to 0

    loop:
        StrCpy $R5 $R2 $R3 $R4 # Create a moving window of the string that is
                               # the size of the length of the search string
        StrCmp $R5 $R1 done    # Is the contents of the window the same as
                               # search string, then done
        StrCmp $R5 "" done     # Is the window empty, then done
        IntOp $R4 $R4 + 1      # Shift the windows one character
        Goto loop              # Repeat

    done:
        StrCpy $R1 $R2 "" $R4
        Pop $R5
        Pop $R4
        Pop $R3
        Pop $R2
        Exch $R1 # $R1=old$R1, stack=[result,...]

FunctionEnd
!macroend
!insertmacro StrStr ""
!insertmacro StrStr "un."


#------------------------------------------------------------------------------
# AddToPath Function
# - Adds item to Path for All Users
# - Overcomes NSIS ReadRegStr limitation of 1024 characters by using Native
#   Windows Commands
#
# Usage:
#   Push "C:\path\to\add"
#   Call AddToPath
#------------------------------------------------------------------------------
!define Environ 'HKLM "SYSTEM\CurrentControlSet\Control\Session Manager\Environment"'
Function AddToPath

    Exch $0 # Path to add
    Push $1 # Current Path
    Push $2 # Results of StrStr / Length of Path + Path to Add
    Push $3 # Handle to Reg / Length of Path
    Push $4 # Result of Registry Call

    # Open a handle to the key in the registry, handle in $3, Error in $4
    System::Call "advapi32::RegOpenKey(i 0x80000002, t'SYSTEM\CurrentControlSet\Control\Session Manager\Environment', *i.r3) i.r4"
    # Make sure registry handle opened successfully (returned 0)
    IntCmp $4 0 0 done done

    # Load the contents of path into $1, Error Code into $4, Path length into $2
    System::Call "advapi32::RegQueryValueEx(i $3, t'PATH', i 0, i 0, t.r1, *i ${NSIS_MAX_STRLEN} r2) i.r4"

    # Close the handle to the registry ($3)
    System::Call "advapi32::RegCloseKey(i $3)"

    # Check for Error Code 234, Path too long for the variable
    IntCmp $4 234 0 +4 +4 # $4 == ERROR_MORE_DATA
        DetailPrint "AddToPath Failed: original length $2 > ${NSIS_MAX_STRLEN}"
        MessageBox MB_OK \
            "You may add C:\clouldops-agent to the %PATH% for convenience when issuing local salt commands from the command line." \
            /SD IDOK
        Goto done

    # If no error, continue
    IntCmp $4 0 +5 # $4 != NO_ERROR
        # Error 2 means the Key was not found
        IntCmp $4 2 +3 # $4 != ERROR_FILE_NOT_FOUND
            DetailPrint "AddToPath: unexpected error code $4"
            Goto done
        StrCpy $1 ""

    # Check if already in PATH
    Push "$1;"          # The string to search
    Push "$0;"          # The string to find
    Call StrStr
    Pop $2              # The result of the search
    StrCmp $2 "" 0 done # String not found, try again with ';' at the end
                        # Otherwise, it's already in the path
    Push "$1;"          # The string to search
    Push "$0\;"         # The string to find
    Call StrStr
    Pop $2              # The result
    StrCmp $2 "" 0 done # String not found, continue (add)
                        # Otherwise, it's already in the path

    # Prevent NSIS string overflow
    StrLen $2 $0        # Length of path to add ($2)
    StrLen $3 $1        # Length of current path ($3)
    IntOp $2 $2 + $3    # Length of current path + path to add ($2)
    IntOp $2 $2 + 2     # Account for the additional ';'
                        # $2 = strlen(dir) + strlen(PATH) + sizeof(";")

    # Make sure the new length isn't over the NSIS_MAX_STRLEN
    IntCmp $2 ${NSIS_MAX_STRLEN} +4 +4 0
        DetailPrint "AddToPath Failed: new length $2 > ${NSIS_MAX_STRLEN}"
        MessageBox MB_OK \
            "You may add C:\salt to the %PATH% for convenience when issuing local salt commands from the command line." \
            /SD IDOK
        Goto done

    # Append dir to PATH
    DetailPrint "Add to PATH: $0"
    StrCpy $2 $1 1 -1       # Copy the last character of the existing path
    StrCmp $2 ";" 0 +2      # Check for trailing ';'
        StrCpy $1 $1 -1     # remove trailing ';'
    StrCmp $1 "" +2         # Make sure Path is not empty
        StrCpy $0 "$1;$0"   # Append new path at the end ($0)

    # We can use the NSIS command here. Only 'ReadRegStr' is affected
    WriteRegExpandStr ${Environ} "PATH" $0

    # Broadcast registry change to open programs
    SendMessage ${HWND_BROADCAST} ${WM_WININICHANGE} 0 "STR:Environment" /TIMEOUT=5000

    done:
        Pop $4
        Pop $3
        Pop $2
        Pop $1
        Pop $0

FunctionEnd


#------------------------------------------------------------------------------
# RemoveFromPath Function
# - Removes item from Path for All Users
# - Overcomes NSIS ReadRegStr limitation of 1024 characters by using Native
#   Windows Commands
#
# Usage:
#   Push "C:\path\to\add"
#   Call un.RemoveFromPath
#------------------------------------------------------------------------------
Function un.RemoveFromPath

    Exch $0
    Push $1
    Push $2
    Push $3
    Push $4
    Push $5
    Push $6

    # Open a handle to the key in the registry, handle in $3, Error in $4
    System::Call "advapi32::RegOpenKey(i 0x80000002, t'SYSTEM\CurrentControlSet\Control\Session Manager\Environment', *i.r3) i.r4"
    # Make sure registry handle opened successfully (returned 0)
    IntCmp $4 0 0 done done

    # Load the contents of path into $1, Error Code into $4, Path length into $2
    System::Call "advapi32::RegQueryValueEx(i $3, t'PATH', i 0, i 0, t.r1, *i ${NSIS_MAX_STRLEN} r2) i.r4"

    # Close the handle to the registry ($3)
    System::Call "advapi32::RegCloseKey(i $3)"

    # Check for Error Code 234, Path too long for the variable
    IntCmp $4 234 0 +4 +4 # $4 == ERROR_MORE_DATA
        DetailPrint "AddToPath: original length $2 > ${NSIS_MAX_STRLEN}"
        Goto done

    # If no error, continue
    IntCmp $4 0 +5 # $4 != NO_ERROR
        # Error 2 means the Key was not found
        IntCmp $4 2 +3 # $4 != ERROR_FILE_NOT_FOUND
            DetailPrint "AddToPath: unexpected error code $4"
            Goto done
        StrCpy $1 ""

    # Ensure there's a trailing ';'
    StrCpy $5 $1 1 -1   # Copy the last character of the path
    StrCmp $5 ";" +2    # Check for trailing ';', if found continue
        StrCpy $1 "$1;" # ensure trailing ';'

    # Check for our directory inside the path
    Push $1             # String to Search
    Push "$0;"          # Dir to Find
    Call un.StrStr
    Pop $2              # The results of the search
    StrCmp $2 "" done   # If results are empty, we're done, otherwise continue

    # Remove our Directory from the Path
    DetailPrint "Remove from PATH: $0"
    StrLen $3 "$0;"       # Get the length of our dir ($3)
    StrLen $4 $2          # Get the length of the return from StrStr ($4)
    StrCpy $5 $1 -$4      # $5 is now the part before the path to remove
    StrCpy $6 $2 "" $3    # $6 is now the part after the path to remove
    StrCpy $3 "$5$6"      # Combine $5 and $6

    # Check for Trailing ';'
    StrCpy $5 $3 1 -1     # Load the last character of the string
    StrCmp $5 ";" 0 +2    # Check for ';'
        StrCpy $3 $3 -1   # remove trailing ';'

    # Write the new path to the registry
    WriteRegExpandStr ${Environ} "PATH" $3

    # Broadcast the change to all open applications
    SendMessage ${HWND_BROADCAST} ${WM_WININICHANGE} 0 "STR:Environment" /TIMEOUT=5000

    done:
        Pop $6
        Pop $5
        Pop $4
        Pop $3
        Pop $2
        Pop $1
        Pop $0

FunctionEnd




Var cfg_line
Var chk_line
Var lst_check
Function updateMinionConfig

    ClearErrors
    FileOpen $0 "$INSTDIR\conf\cloudops.yaml" "r"       # open target file for reading
    GetTempFileName $R0                                 # get new temp file name
    FileOpen $1 $R0 "w"                                 # open temp file for writing

    StrCpy $ConfigWriteMaster 1                         # write the master config value
    StrCpy $ConfigWriteMinion 1                         # write the minion config value

    loop:                                               # loop through each line
        FileRead $0 $cfg_line                           # read line from target file
        IfErrors done                                   # end if errors are encountered (end of line)

        loop_after_read:
        StrCpy $lst_check 0                             # list check not performed

        ${If} $ConfigWriteMaster == 1                   # if we need to write master config

            ${StrLoc} $3 $cfg_line "api_key:" ">"        # where is 'master:' in this line
            ${If} $3 == 0                               # is it in the first...
            ${OrIf} $3 == 1                             # or second position (account for comments)

                ${Explode} $9 "," $ApiKey_state     # Split the hostname on commas, $9 is the number of items found
                ${If} $9 == 1                           # 1 means only a single master was passed
                    StrCpy $cfg_line "api_key: $ApiKey_State$\r$\n"  # write the master
                ${Else}                                 # make a multi-master entry
                    StrCpy $cfg_line "api_key:"          # make the first line "master:"

                    loop_explode:                       # start a loop to go through the list in the config
                    pop $8                              # pop the next item off the stack
                    ${Trim} $8 $8                       # trim any whitespace
                    StrCpy $cfg_line "$cfg_line$\r$\n  - $8"  # add it to the master variable ($2)
                    IntOp $9 $9 - 1                     # decrement the list count
                    ${If} $9 >= 1                       # if it's not 0
                        Goto loop_explode               # do it again
                    ${EndIf}                            # close if statement
                    StrCpy $cfg_line "$cfg_line$\r$\n"  # Make sure there's a new line at the end

                    # Remove remaining items in list
                    ${While} $lst_check == 0            # while list item found
                        FileRead $0 $chk_line           # read line from target file
                        IfErrors done                   # end if errors are encountered (end of line)
                        ${StrLoc} $3 $chk_line "  - " ">"  # where is 'master:' in this line
                        ${If} $3 == ""                  # is it in the first...
                            StrCpy $lst_check 1         # list check performed and finished
                        ${EndIf}
                    ${EndWhile}

                ${EndIf}                                # close if statement

                StrCpy $ConfigWriteMaster 0             # master value written to config

            ${EndIf}                                    # close if statement
        ${EndIf}                                        # close if statement

        ${If} $ConfigWriteMinion == 1                   # if we need to write minion config
            ${StrLoc} $3 $cfg_line "endpoint:" ">"            # where is 'id:' in this line
            ${If} $3 == 0                               # is it in the first...
            ${OrIf} $3 == 1                             # or the second position (account for comments)
                StrCpy $cfg_line "endpoint: $Endpoint_State$\r$\n"  # write the minion config setting
                StrCpy $ConfigWriteMinion 0             # minion value written to config
            ${EndIf}                                    # close if statement
        ${EndIf}                                        # close if statement

        FileWrite $1 $cfg_line                          # write changed or unchanged line to temp file

    ${If} $lst_check == 1                               # master not written to the config
        StrCpy $cfg_line $chk_line
        Goto loop_after_read                            # A loop was performed, skip the next read
    ${EndIf}                                            # close if statement

    Goto loop                                           # check the next line in the config file

    done:
    ClearErrors
    # Does master config still need to be written
    ${If} $ConfigWriteMaster == 1                       # master not written to the config

        ${Explode} $9 "," $ApiKey_State             # split the hostname on commas, $9 is the number of items found
        ${If} $9 == 1                                   # 1 means only a single master was passed
            StrCpy $cfg_line "api_key: $ApiKey_State"  # write the master
        ${Else}                                         # make a multi-master entry
            StrCpy $cfg_line "api_key:"                  # make the first line "master:"

            loop_explode_2:                             # start a loop to go through the list in the config
            pop $8                                      # pop the next item off the stack
            ${Trim} $8 $8                               # trim any whitespace
            StrCpy $cfg_line "$cfg_line$\r$\n  - $8"    # add it to the master variable ($2)
            IntOp $9 $9 - 1                             # decrement the list count
            ${If} $9 >= 1                               # if it's not 0
                Goto loop_explode_2                     # do it again
            ${EndIf}                                    # close if statement
        ${EndIf}                                        # close if statement
        FileWrite $1 $cfg_line                          # write changed or unchanged line to temp file

    ${EndIf}                                            # close if statement

    ${If} $ConfigWriteMinion == 1                       # minion ID not written to the config
        StrCpy $cfg_line "$\r$\nendpoint: $Endpoint_State"  # write the minion config setting
        FileWrite $1 $cfg_line                          # write changed or unchanged line to temp file
    ${EndIf}                                            # close if statement

    FileClose $0                                        # close target file
    FileClose $1                                        # close temp file
    Delete "$INSTDIR\conf\cloudops.yaml"                       # delete target file
    CopyFiles /SILENT $R0 "$INSTDIR\conf\cloudops.yaml"        # copy temp file to target file
    Delete $R0                                          # delete temp file

FunctionEnd


Function parseCommandLineSwitches

    # Load the parameters
    ${GetParameters} $R0

    # Display Help
    ClearErrors
    ${GetOptions} $R0 "/?" $R1
    IfErrors display_help_not_found

        System::Call 'kernel32::GetStdHandle(i -11)i.r0'
        System::Call 'kernel32::AttachConsole(i -1)i.r1'
        ${If} $0 = 0
        ${OrIf} $1 = 0
            System::Call 'kernel32::AllocConsole()'
            System::Call 'kernel32::GetStdHandle(i -11)i.r0'
        ${EndIf}
        FileWrite $0 "$\n"
        FileWrite $0 "$\n"
        FileWrite $0 "Help for CloudOps Agent installation$\n"
        FileWrite $0 "===============================================================================$\n"
        FileWrite $0 "$\n"
        FileWrite $0 "/endpoint=$\t$\tA string value to set the endpoint. $\n"
        FileWrite $0 "/apikey=$\t$\tA string value to set the apiKey. $\n"
        FileWrite $0 "-------------------------------------------------------------------------------$\n"
        FileWrite $0 "$\n"
        FileWrite $0 "Examples:$\n"
        FileWrite $0 "$\n"
        FileWrite $0 "${OutFile} /S$\n"
        FileWrite $0 "$\n"
        FileWrite $0 "${OutFile} /S /endpoint=endpoint /apiKey=apiKey $\n"
        FileWrite $0 "$\n"
        FileWrite $0 "===============================================================================$\n"
        FileWrite $0 "$\n"
        System::Free $0
        System::Free $1
        System::Call 'kernel32::FreeConsole()'

        # Give the user back the prompt
        !define VK_RETURN 0x0D ; Enter Key
        !define KEYEVENTF_EXTENDEDKEY 0x0001
        !define KEYEVENTF_KEYUP 0x0002
        System::Call "user32::keybd_event(i${VK_RETURN}, i0x45, i${KEYEVENTF_EXTENDEDKEY}|0, i0)"
        System::Call "user32::keybd_event(i${VK_RETURN}, i0x45, i${KEYEVENTF_EXTENDEDKEY}|${KEYEVENTF_KEYUP}, i0)"
        Abort

    display_help_not_found:

    # Check for start-minion switches
    # /start-service is to be deprecated, so we must check for both
    ${GetOptions} $R0 "/start-service=" $R1
    ${GetOptions} $R0 "/start-minion=" $R2

    # Service: Start Salt Minion
    ${IfNot} $R2 == ""
        # If start-minion was passed something, then set it
        StrCpy $StartMinion $R2
    ${ElseIfNot} $R1 == ""
        # If start-service was passed something, then set StartMinion to that
        StrCpy $StartMinion $R1
    ${Else}
        # Otherwise default to 1
        StrCpy $StartMinion 1
    ${EndIf}

    # Minion Config: Master IP/Name
    # If setting master, we don't want to use existing config
    ${GetOptions} $R0 "/endpoint=" $R1
    ${IfNot} $R1 == ""
        StrCpy $Endpoint_State $R1
    ${ElseIf} $Endpoint_State == ""
        StrCpy $Endpoint_State ""
    ${EndIf}

    # Minion Config: Minion ID
    # If setting minion id, we don't want to use existing config
    ${GetOptions} $R0 "/apiKey=" $R1
    ${IfNot} $R1 == ""
        StrCpy $ApiKey_State $R1
    ${ElseIf} $ApiKey_State == ""
        StrCpy $ApiKey_State ""
    ${EndIf}


FunctionEnd

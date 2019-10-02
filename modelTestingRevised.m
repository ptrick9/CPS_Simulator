

%Create Nodes
numNodes = 1000

%Create Time Steps
t = [0:1:20000]';

%Create All S0-S2 parameters
S0 = randn(1,numNodes)*0.005+0.33;
S1 = randn(1,numNodes)*0.005+0.33;
S2 = randn(1,numNodes)*0.005+0.33;

%Create all E0-E2 parameters
E0 = randn(1,numNodes)*0.02.*S0;
E1 = randn(1,numNodes)*0.02.*S1;
E2 = randn(1,numNodes)*0.02.*S2;

%Tau1-2 as defined in original document
Tau1 = 3500;
Tau2 = 9000;

%ET1-2 values 
ET1 = Tau1*randn(1,numNodes)*0.05;
ET2 = Tau1*randn(1,numNodes)*0.05;


%Sensitivity for each node over time
S = (S0+E0)+(S1+E1).*exp(-t./(Tau1+ET1))+(S2+E2).*exp(-t./(Tau2+ET2));

%Plot to see range and falloff
plot(S)


%Calculate average halfway point for dropoff
%Some nodes parameters don't have them drop past 1/2 sensitivity so set
%initial values to number of Timesteps so that average is 'correct'
half = ones(1,numNodes)*10000;
long = 0;
for i=1:1:numNodes
    node = S(:,i)';
    a = node <= node(1)/2;
    ind = find(a, 1, 'first');
    if isempty(ind) 
        half(i) = 0;
        long = long + 1;
    else
        half(i) = ind;
    end

end

%report average halfway dropoff time
mean(half)
